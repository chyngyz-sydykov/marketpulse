package grpc

import (
	"context"
	"log"
	"net"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	pb "github.com/chyngyz-sydykov/marketpulse/proto/marketpulse"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	pb.MarketPulseServer
	MarketService     *marketdata.MarketDataService
	IndicatorService  *indicator.IndicatorService
	MarketDataHandler *MarketDataHandler
	IndicatorHandler  *IndicatorHandler
}

func NewGrpcService(MarketService *marketdata.MarketDataService, IndicatorService *indicator.IndicatorService) *GrpcServer {
	MarketDataHandler := NewMarketDataHandler(*MarketService)
	IndicatorHandler := NewIndicatorHandler(*IndicatorService)

	return &GrpcServer{
		MarketService:     MarketService,
		IndicatorService:  IndicatorService,
		MarketDataHandler: MarketDataHandler,
		IndicatorHandler:  IndicatorHandler,
	}
}
func (server *GrpcServer) GetOHLC(ctx context.Context, request *pb.OHLCRequest) (*pb.OHLCResponse, error) {
	return server.MarketDataHandler.GetOHLC(ctx, request)
}

func (server *GrpcServer) GetIndicators(ctx context.Context, request *pb.IndicatorRequest) (*pb.IndicatorResponse, error) {
	return server.IndicatorHandler.GetIndicators(ctx, request)
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
