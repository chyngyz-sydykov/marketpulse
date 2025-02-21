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
	_, err := suite.db.Exec("DELETE FROM data_btc; DELETE FROM indicator_btc;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
}

func (suite *GrpcTestSuite) TestGrpcMarketDataEndpointShouldReturnCorrectData() {
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

func (suite *GrpcTestSuite) TestGrpcIndicatorEndpointShouldReturnCorrectData() {
	// Arrange: Insert mock indicator data
	records := []*dto.DataDto{
		{Symbol: "BTCUSDT", Timeframe: "4h", Timestamp: time.Date(2020, 1, 1, 4, 0, 0, 0, time.Now().Location()), Open: 110, High: 120, Low: 100, Close: 105, Volume: 111, Trend: -0.005, IsComplete: true},
		{Symbol: "BTCUSDT", Timeframe: "4h", Timestamp: time.Date(2020, 1, 1, 8, 0, 0, 0, time.Now().Location()), Open: 211, High: 221, Low: 201, Close: 206, Volume: 223, IsComplete: true},
	}

	suite.marketDataService.UpsertBatchData("btc", records)

	_, err := suite.db.Exec(`INSERT INTO indicator_btc_4h (timeframe, timestamp,data_timestamp, sma, ema, std_dev, lower_bollinger, upper_bollinger, rsi, volatility, macd, macd_signal) VALUES 
	('4h', '2020-01-01 04:00:00', '2020-01-01 04:00:00', 105, 107, 2.5, 100, 110, 45, 0.03, 0.5, 0.3),
	('4h', '2020-01-01 08:00:00', '2020-01-01 08:00:00', 205, 207, 3.1, 200, 210, 50, 0.02, 0.6, 0.4)`)
	if err != nil {
		log.Fatalf("Failed to insert mock indicator data: %v", err)
	}
	listener, server := createInMemoryGrpcServer(suite)
	defer server.Stop()

	// Create gRPC client
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewMarketPulseClient(conn)

	// ✅ Define test cases
	testCases := []struct {
		name       string
		request    *pb.IndicatorRequest
		expectLen  int
		expectData *pb.IndicatorResponse
		expectErr  string
	}{
		{
			name:      "Fetch 4H indicator data with end time filter",
			request:   &pb.IndicatorRequest{Currency: "btc", Timeframe: "4h", EndTime: timestamppb.New(time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC))},
			expectLen: 1,
			expectData: &pb.IndicatorResponse{Indicators: []*pb.IndicatorData{
				{Timeframe: "4h", Sma: 105, Ema: 107, StdDev: 2.5, LowerBollinger: 100, UpperBollinger: 110, Rsi: 45, Volatility: 0.03, Macd: 0.5, MacdSignal: 0.3},
			}},
		},
		{
			name:      "Fetch 4H indicator data with start time filter (empty response)",
			request:   &pb.IndicatorRequest{Currency: "btc", Timeframe: "4h", StartTime: timestamppb.New(time.Date(2022, 1, 1, 6, 0, 0, 0, time.UTC))},
			expectLen: 0,
		},
		{
			name:      "Fetch all 4H indicator data",
			request:   &pb.IndicatorRequest{Currency: "btc", Timeframe: "4h"},
			expectLen: 2,
		},
		{
			name:      "Fetch 1D indicator data (empty response)",
			request:   &pb.IndicatorRequest{Currency: "btc", Timeframe: "1d"},
			expectLen: 0,
		},
		{
			name:      "Invalid currency",
			request:   &pb.IndicatorRequest{Currency: "unknown currency", Timeframe: "1h"},
			expectErr: "resource value(s) is invalid",
		},
	}

	// ✅ Execute test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			response, err := client.GetIndicators(context.Background(), tc.request)
			// ✅ Handle expected errors
			if tc.expectErr != "" {
				assert.Error(suite.T(), err)
				assert.ErrorContains(suite.T(), err, tc.expectErr)
				return
			}

			// ✅ Validate response
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tc.expectLen, len(response.Indicators))

			// ✅ Validate indicator data
			if tc.expectData != nil {
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Timeframe, response.Indicators[0].Timeframe)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Sma, response.Indicators[0].Sma)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Ema, response.Indicators[0].Ema)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].StdDev, response.Indicators[0].StdDev)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].LowerBollinger, response.Indicators[0].LowerBollinger)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].UpperBollinger, response.Indicators[0].UpperBollinger)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Rsi, response.Indicators[0].Rsi)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Volatility, response.Indicators[0].Volatility)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].Macd, response.Indicators[0].Macd)
				assert.Equal(suite.T(), tc.expectData.Indicators[0].MacdSignal, response.Indicators[0].MacdSignal)
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
