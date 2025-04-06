package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type API struct {
	server *http.Server
}

func main() {
	r := gin.Default()

	httpPort := getEnv("HTTP_PORT", "8003")
	inventoryServiceURL := getEnv("INVENTORY_SERVICE", "http://0.0.0.0:8001")
	orderServiceURL := getEnv("ORDER_SERVICE", "http://0.0.0.0:8002")
	httpMode := getEnv("GIN_MODE", "release")

	gin.SetMode(httpMode)

	inventoryProxy := createProxy(inventoryServiceURL)
	orderProxy := createProxy(orderServiceURL)

	r.Any("/api/v1/products/*path", proxyHandler(inventoryProxy))
	r.Any("/api/v1/categories/*path", proxyHandler(inventoryProxy))
	r.Any("/api/v1/orders/*path", proxyHandler(orderProxy))

	server := &http.Server{
		Addr:    "0.0.0.0:" + httpPort,
		Handler: r,
	}

	api := &API{server: server}

	go func() {
		log.Printf("Gateway running on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	api.waitForShutdown()
}

func (a *API) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("Shutdown signal received", "signal:", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully...")

	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("HTTP server stopped successfully")
}

func createProxy(target string) *httputil.ReverseProxy {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Host = targetURL.Host
	}

	return proxy
}

func proxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		originalPath := c.Request.URL.Path

		cleanedPath := strings.ReplaceAll(originalPath, "//", "/")
		c.Request.URL.Path = cleanedPath

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
