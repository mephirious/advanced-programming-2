package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mephirious/advanced-programming-2/order-service/config"
	"github.com/mephirious/advanced-programming-2/order-service/internal/adapter/http/service"
	"github.com/mephirious/advanced-programming-2/order-service/internal/adapter/http/service/gateway"
	"github.com/mephirious/advanced-programming-2/order-service/internal/repository"
	"github.com/mephirious/advanced-programming-2/order-service/internal/usecase"
	"github.com/mephirious/advanced-programming-2/order-service/pkg/mongo"
)

const serviceName = "order-service"

type App struct {
	httpServer *service.API
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service", serviceName)

	log.Println("connecting to mongo", "database", cfg.Mongo.Database)
	mongoDB, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	gatewayClient := gateway.NewHTTPGateway(cfg.Server.GatewayService)
	orderRepository := repository.NewOrderRepository(mongoDB.Connection)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, gatewayClient)

	httpServer := service.New(cfg.Server, orderUseCase)

	app := &App{
		httpServer: httpServer,
	}

	return app, nil
}

func (a *App) Close() {
	err := a.httpServer.Stop()
	if err != nil {
		log.Println("failed to shutdown http service", err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	a.httpServer.Run(errCh)

	log.Printf("service %v started", serviceName)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...", s)

		a.Close()
		log.Println("graceful shutdown completed!")
	}

	return nil
}
