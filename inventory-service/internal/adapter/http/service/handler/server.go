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
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string
}

func New(cfg config.Server) *API {
	gin.SetMode(cfg.HTTPServer.Mode)

	server := gin.New()
	server.Use(gin.Recovery())

	api := &API{
		server: server,
		cfg:    cfg.HTTPServer,
		addr:   fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
	}

	api.setupRoutes()
	return api
}

func (a *API) setupRoutes() {
	//v1 := a.server.Group("/api/v1")
	{
		//
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
	log.Println("Shutdown signal received:", sig.String())

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")
	log.Println("HTTP server stopped successfully")

	return nil
}
