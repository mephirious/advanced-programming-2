package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mephirious/advanced-programming-2/order-service/config"
	"github.com/mephirious/advanced-programming-2/order-service/internal/adapter/grpc/service"
	producer "github.com/mephirious/advanced-programming-2/order-service/internal/adapter/nats"
	"github.com/mephirious/advanced-programming-2/order-service/internal/repository"
	"github.com/mephirious/advanced-programming-2/order-service/internal/usecase"
	"github.com/mephirious/advanced-programming-2/order-service/pkg/mongo"
	"github.com/mephirious/advanced-programming-2/order-service/pkg/nats"
)

const serviceName = "order-service"

type App struct {
	grpcServer *service.GRPCServer
	natsClient *nats.Client
	orderProd  *producer.OrderEventProducer
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %s", serviceName)

	mongoDB, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	natsClient, err := nats.NewClient(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("nats.NewClient: %w", err)
	}
	orderProducer := producer.NewOrderEventProducer(natsClient, "inventory.events")

	orderRepo := repository.NewOrderRepository(mongoDB.Connection)
	orderUC := usecase.NewOrderUseCase(orderRepo, *orderProducer)

	grpcServer, err := service.NewGRPCServer(*cfg, orderUC)
	if err != nil {
		return nil, err
	}

	return &App{
		grpcServer: grpcServer,
		natsClient: natsClient,
		orderProd:  orderProducer,
	}, nil
}

func (a *App) Close() {
	a.grpcServer.Stop()
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- a.grpcServer.Run()
	}()

	log.Printf("%s started", serviceName)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-shutdownCh:
		log.Printf("received signal: %v. Gracefully shutting down...", sig)
		a.Close()
		log.Println("shutdown completed")
	}

	return nil
}
