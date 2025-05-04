package service

import (
	"fmt"
	"log"
	"net"

	"github.com/mephirious/advanced-programming-2/inventory-service/config"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/grpc/service/handler"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
	pb "github.com/mephirious/advanced-programming-2/inventory-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	cfg      config.GRPCServer
	server   *grpc.Server
	addr     string
	listener net.Listener
}

func NewGRPCServer(cfg config.Config, productUC usecase.ProductUseCase, categoryUC usecase.CategoryUseCase) (*GRPCServer, error) {
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.GRPCServer.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s := grpc.NewServer()
	handler := handler.NewInventoryHandler(productUC, categoryUC)

	pb.RegisterInventoryServiceServer(s, handler)

	reflection.Register(s)

	return &GRPCServer{
		cfg:      cfg.Server.GRPCServer,
		server:   s,
		addr:     addr,
		listener: lis,
	}, nil
}

func (s *GRPCServer) Run() error {
	log.Printf("gRPC server running on %s", s.addr)
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	log.Println("gRPC server stopped gracefully")
}
