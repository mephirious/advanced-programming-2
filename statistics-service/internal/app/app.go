package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mephirious/advanced-programming-2/statistics-service/config"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/adapter/grpc"
	handler "github.com/mephirious/advanced-programming-2/statistics-service/internal/adapter/nats/handler"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/repository"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/usecase"
	"github.com/mephirious/advanced-programming-2/statistics-service/pkg/mongo"
	nats "github.com/mephirious/advanced-programming-2/statistics-service/pkg/nats"
)

const serviceName = "statistics-service"

type App struct {
	grpcServer  *grpc.Server
	natsHandler *handler.NATSHandler
	mongoDB     *mongo.DB
	natsConn    *nats.Client
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service", serviceName)

	log.Println("connecting to mongo", "database", cfg.Mongo.Database)
	mongoDB, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	log.Println("connecting to NATS", "url", cfg.NATS.URL)
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %w", err)
	}

	repo := repository.NewMongoStatsRepository(mongoDB.Connection)
	uc := usecase.NewStatsUseCase(repo)

	grpcServer, err := grpc.NewServer(cfg.Server.GRPCServer.Port, uc)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC server: %w", err)
	}

	natsHandler := handler.NewNATSHandler(uc, nc.Conn)

	return &App{
		grpcServer:  grpcServer,
		natsHandler: natsHandler,
		mongoDB:     mongoDB,
		natsConn:    nc,
	}, nil
}

func (a *App) Close() {
	a.grpcServer.Stop()
	a.natsConn.Close()
}

func (a *App) Run() error {
	go func() {
		if err := a.natsHandler.Start(); err != nil {
			log.Fatalf("NATS handler failed: %v", err)
		}
	}()

	go func() {
		if err := a.grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	log.Printf("service %v started", serviceName)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...", s)
		a.Close()
		log.Println("graceful shutdown completed!")
	}

	return nil
}
