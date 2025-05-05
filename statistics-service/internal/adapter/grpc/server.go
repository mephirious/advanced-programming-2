package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/adapter/grpc/handler"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/usecase"
	pb "github.com/mephirious/advanced-programming-2/statistics-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewServer(port int, uc usecase.StatsUseCase) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	handler := handler.NewGRPCHandler(uc)
	pb.RegisterStatisticsServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return &Server{
		server:   grpcServer,
		listener: lis,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("Starting gRPC server on %s", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

func (s *Server) Stop() {
	log.Println("Stopping gRPC server gracefully...")
	s.server.GracefulStop()
}
