package grpc

import (
	"context"
	"log"
	"net"

	"github.com/chyngyz-sydykov/marketpulse/config"
	pb "github.com/chyngyz-sydykov/marketpulse/proto/marketpulse"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	pb.MarketPulseServer
}

func (s *GrpcServer) GetOHLC(ctx context.Context, in *pb.OHLCRequest) (*pb.OHLCResponse, error) {
	// Implement your service method logic here
	return &pb.OHLCResponse{}, nil
}
func (s *GrpcServer) GetIndicators(ctx context.Context, in *pb.IndicatorRequest) (*pb.IndicatorResponse, error) {
	// Implement your service method logic here
	return &pb.IndicatorResponse{}, nil
}
func StartGRPCServer() {
	cfg := config.LoadConfig()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterMarketPulseServer(grpcServer, &GrpcServer{})

	listener, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GrpcPort, err)
	}

	log.Println("Rating service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
