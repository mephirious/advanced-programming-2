package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mephirious/advanced-programming-2/inventory-service/config"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/cache"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/grpc/service"
	producer "github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/nats"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/repository"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/mongo"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/nats"
)

const serviceName = "inventory-service"

type App struct {
	grpcServer    *service.GRPCServer
	natsClient    *nats.Client
	inventoryProd *producer.InventoryEventProducer
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service", serviceName)

	log.Println("connecting to mongo", "database", cfg.Mongo.Database)
	mongoDB, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	natsClient, err := nats.NewClient(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("nats.NewClient: %w", err)
	}
	inventoryProducer := producer.NewInventoryEventProducer(natsClient, "inventory.events")
	productCache := cache.NewProductCache()

	productRepository := repository.NewProductRepository(mongoDB.Connection)
	productUseCase := usecase.NewProductUseCase(productRepository, inventoryProducer, productCache)
	cache.StartCacheRefresher(productCache)

	categoryRepository := repository.NewCategoryRepository(mongoDB.Connection)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)

	grpcServer, err := service.NewGRPCServer(*cfg, productUseCase, categoryUseCase)
	if err != nil {
		return nil, err
	}

	return &App{
		grpcServer:    grpcServer,
		natsClient:    natsClient,
		inventoryProd: inventoryProducer,
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
