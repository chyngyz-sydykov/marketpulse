package grpc

import (
	"context"
	"log"
	"net"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	pb "github.com/chyngyz-sydykov/marketpulse/proto/marketpulse"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pb.MarketPulseServer
	MarketService    *marketdata.MarketDataService
	IndicatorService *indicator.IndicatorService
}

func NewGrpcService(MarketService *marketdata.MarketDataService, IndicatorService *indicator.IndicatorService) *GrpcServer {
	return &GrpcServer{
		MarketService:    MarketService,
		IndicatorService: IndicatorService,
	}
}

func (server *GrpcServer) GetOHLC(ctx context.Context, request *pb.OHLCRequest) (*pb.OHLCResponse, error) {
	ohlcRequestDto := server.mapProtoRequestToDTO(request)
	records, err := server.MarketService.GetRecordsByRequest(ohlcRequestDto)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return nil, status.Error(codes.InvalidArgument, "resource value(s) is invalid")
	}

	ohlcResponse := server.mapDataDtoToOHLCResponse(records)
	return ohlcResponse, nil
}

func (server *GrpcServer) GetIndicators(ctx context.Context, in *pb.IndicatorRequest) (*pb.IndicatorResponse, error) {
	// Implement your service method logic here
	return &pb.IndicatorResponse{}, nil
}
func (server *GrpcServer) mapProtoRequestToDTO(req *pb.OHLCRequest) dto.OHLCRequestDto {
	dto := dto.OHLCRequestDto{
		Currency:  req.Currency,
		Timeframe: req.Timeframe,
	}

	// Handle optional timestamp fields
	if req.StartTime != nil {
		startTime := req.StartTime.AsTime()
		dto.StartTime = &startTime
	}

	if req.EndTime != nil {
		endTime := req.EndTime.AsTime()
		dto.EndTime = &endTime
	}

	// Handle optional fields with default values
	if req.Limit > 0 {
		dto.Limit = req.Limit
	} else {
		dto.Limit = int32(config.DEFAULT_LIMIT)
	}

	if req.SortField != "" && slices.Contains([]string{"timestamp", "volume"}, req.SortField) {
		dto.SortField = req.SortField
	} else {
		dto.SortField = config.DEFAULT_DATA_REQUEST_SORT_FIELD
	}
	if req.SortOrder != "" && (req.SortOrder == "ASC" || req.SortOrder == "DESC") {
		dto.SortOrder = req.SortOrder
	} else {
		dto.SortOrder = config.DEFAULT_DATA_REQUEST_SORT_ORDER
	}

	return dto
}
func (server *GrpcServer) mapDataDtoToOHLCResponse(records []dto.DataDto) *pb.OHLCResponse {
	ohlcRecords := make([]*pb.OHLCData, len(records))
	for i, record := range records {
		ohlcRecords[i] = &pb.OHLCData{
			Timestamp:  timestamppb.New(record.Timestamp),
			Timeframe:  record.Timeframe,
			Currency:   record.Symbol,
			Open:       record.Open,
			High:       record.High,
			Low:        record.Low,
			Close:      record.Close,
			Volume:     record.Volume,
			Trend:      record.Trend,
			IsComplete: record.IsComplete,
		}
	}
	return &pb.OHLCResponse{
		Records: ohlcRecords,
	}
}

func StartGRPCServer(GrpcServer *GrpcServer) {
	cfg := config.LoadConfig()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterMarketPulseServer(grpcServer, GrpcServer)

	listener, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GrpcPort, err)
	}

	log.Println("Rating service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
