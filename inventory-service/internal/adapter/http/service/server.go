package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mephirious/advanced-programming-2/inventory-service/config"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/http/service/handler"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string

	productHandler  *handler.ProductHandler
	categoryHandler *handler.CategoryHandler
	discountHandler *handler.DiscountHandler
}

func New(cfg config.Server, productUsecase usecase.ProductUseCase, categoryUseCase usecase.CategoryUseCase, discountUseCase usecase.DiscountUseCase) *API {
	gin.SetMode(cfg.HTTPServer.Mode)
	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(gin.Logger())

	productHandler := handler.NewProductHandler(productUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	discountHandler := handler.NewDiscountHandler(discountUseCase)

	api := &API{
		server:          server,
		cfg:             cfg.HTTPServer,
		addr:            fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
		discountHandler: discountHandler,
	}

	api.setupRoutes()

	return api
}

func (a *API) setupRoutes() {
	v1 := a.server.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("/", a.productHandler.CreateProduct)
			products.GET("/", a.productHandler.GetAllProducts)
			products.GET("/:id", a.productHandler.GetProductByID)
			products.PATCH("/:id", a.productHandler.UpdateProduct)
			products.DELETE("/:id", a.productHandler.DeleteProduct)
		}

		category := v1.Group("/categories")
		{
			category.POST("/", a.categoryHandler.CreateCategory)
			category.GET("/", a.categoryHandler.GetAllCategories)
			category.GET("/:id", a.categoryHandler.GetCategoryByID)
			category.PATCH("/:id", a.categoryHandler.UpdateCategory)
			category.DELETE("/:id", a.categoryHandler.DeleteCategory)
		}

		discount := v1.Group("/discounts")
		{
			discount.GET("/:id", a.discountHandler.GetAllProductsWithPromotion)
			discount.DELETE("/:id", a.discountHandler.DeleteDiscount)
			discount.POST("/", a.discountHandler.CreateDiscount)
		}
	}
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server starting on: %v", a.addr)

		if err := a.server.Run(a.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) Stop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("Shutdown signal received", "signal:", sig.String())

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")

	log.Println("HTTP server stopped successfully")

	return nil
}
