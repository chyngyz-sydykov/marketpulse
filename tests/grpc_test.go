package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/database"
	mygrpc "github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/grpc"
	pb "github.com/chyngyz-sydykov/marketpulse/proto/marketpulse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Test Suite Struct
type GrpcTestSuite struct {
	suite.Suite
	db                *sql.DB
	marketDataService *marketdata.MarketDataService
	indicatorService  *indicator.IndicatorService
	redisMock         *MockRedisService
}

func (suite *GrpcTestSuite) SetupSuite() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	mockRedis := &MockRedisService{}
	mockRedis.On("PublishEvent", mock.Anything, mock.Anything, "MarketPulse").Return(nil)

	suite.marketDataService = marketdata.NewMarketDataService(mockRedis)
	suite.indicatorService = indicator.NewIndicatorService(mockRedis)
	suite.db = database.DB
	suite.redisMock = mockRedis
}

// Clear DB Before Each Test
func (suite *GrpcTestSuite) SetupTest() {
	_, err := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
}

// ✅ Test Case 1: Store 1H Record
func (suite *GrpcTestSuite) TestGrpcShouldReturnCorrectData() {
	records := []*dto.DataDto{
		{Symbol: "BTCUSDT", Timeframe: "1h", Timestamp: time.Date(2020, 1, 1, 1, 0, 0, 0, time.Now().Location()), Open: 110, High: 120, Low: 100, Close: 105, Volume: 111, Trend: -0.005, IsComplete: true},
		{Symbol: "BTCUSDT", Timeframe: "1h", Timestamp: time.Date(2022, 1, 1, 1, 0, 0, 0, time.Now().Location()), Open: 211, High: 221, Low: 201, Close: 206, Volume: 223, IsComplete: true},
	}

	suite.marketDataService.UpsertBatchData("btc", records)

	record := &dto.DataDto{
		Symbol:     "BTCUSDT",
		Timeframe:  "4h",
		Timestamp:  time.Date(2020, 1, 1, 2, 0, 0, 0, time.Now().Location()),
		Open:       310,
		High:       320,
		Low:        300,
		Close:      305,
		Volume:     311,
		IsComplete: false,
	}

	suite.marketDataService.StoreData("btc", record)

	listener, server := createInMemoryGrpcServer(suite)
	defer server.Stop()

	//Create a gRPC client connected to the server
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewMarketPulseClient(conn)

	testCases := []struct {
		name       string
		request    *pb.OHLCRequest
		expectLen  int
		expectData *pb.OHLCResponse
		expectErr  string
	}{
		{
			name:      "Fetch 1H data",
			request:   &pb.OHLCRequest{Currency: "btc", Timeframe: "1h", EndTime: timestamppb.New(time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC))},
			expectLen: 1,
			expectData: &pb.OHLCResponse{Records: []*pb.OHLCData{
				{Currency: "BTCUSDT", Timeframe: "1h", Open: 110, High: 120, Low: 100, Close: 105, Volume: 111, IsComplete: true},
			}},
		},
		{
			name:      "Fetch 4H data",
			request:   &pb.OHLCRequest{Currency: "btc", Timeframe: "4h"},
			expectLen: 1,
			expectData: &pb.OHLCResponse{Records: []*pb.OHLCData{
				{Currency: "BTCUSDT", Timeframe: "4h", Open: 310, High: 320, Low: 300, Close: 305, Volume: 311, IsComplete: false},
			}},
		},
		{
			name:      "Fetch 1D data (empty response)",
			request:   &pb.OHLCRequest{Currency: "btc", Timeframe: "1d"},
			expectLen: 0,
		},
		{
			name:      "Fetch all 1H data",
			request:   &pb.OHLCRequest{Currency: "btc", Timeframe: "1h"},
			expectLen: 2,
		},
		{
			name:      "Invalid currency",
			request:   &pb.OHLCRequest{Currency: "unknown currency", Timeframe: "1h"},
			expectErr: "resource value(s) is invalid",
		},
	}

	// ✅ Loop through test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			response, err := client.GetOHLC(context.Background(), tc.request)

			// ✅ Check expected error
			if tc.expectErr != "" {
				assert.Error(suite.T(), err)
				assert.ErrorContains(suite.T(), err, tc.expectErr)
				return
			}

			// ✅ Check response
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tc.expectLen, len(response.Records))

			// ✅ Validate response data
			if tc.expectData != nil {
				assert.Equal(suite.T(), tc.expectData.Records[0].Currency, response.Records[0].Currency)
				assert.Equal(suite.T(), tc.expectData.Records[0].Timeframe, response.Records[0].Timeframe)
				assert.Equal(suite.T(), tc.expectData.Records[0].Open, response.Records[0].Open)
				assert.Equal(suite.T(), tc.expectData.Records[0].High, response.Records[0].High)
				assert.Equal(suite.T(), tc.expectData.Records[0].Low, response.Records[0].Low)
				assert.Equal(suite.T(), tc.expectData.Records[0].Close, response.Records[0].Close)
				assert.Equal(suite.T(), tc.expectData.Records[0].Volume, response.Records[0].Volume)
				assert.Equal(suite.T(), tc.expectData.Records[0].IsComplete, response.Records[0].IsComplete)
			}
		})
	}

}

// ✅ Run the Test Suite
func TestGrpcServer(t *testing.T) {
	suite.Run(t, new(GrpcTestSuite))
}

func createInMemoryGrpcServer(suite *GrpcTestSuite) (net.Listener, *grpc.Server) {

	server := grpc.NewServer()
	pb.RegisterMarketPulseServer(server, mygrpc.NewGrpcService(suite.marketDataService, suite.indicatorService))

	listener, err := net.Listen("tcp", "localhost:11111")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	return listener, server
}
