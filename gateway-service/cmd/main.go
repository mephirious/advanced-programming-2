package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventorypb "github.com/mephirious/advanced-programming-2/gateway-service/proto/inventory"
	orderpb "github.com/mephirious/advanced-programming-2/gateway-service/proto/order"
	statpb "github.com/mephirious/advanced-programming-2/gateway-service/proto/statistics"
)

func main() {
	r := gin.Default()
	gin.SetMode(getEnv("GIN_MODE", "release"))

	orderConn, err := grpc.NewClient(getEnv("ORDER_SERVICE_GRPC", "localhost:8002"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	inventoryConn, err := grpc.NewClient(getEnv("INVENTORY_SERVICE_GRPC", "localhost:8001"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := inventorypb.NewInventoryServiceClient(inventoryConn)

	statConn, err := grpc.NewClient(getEnv("STATISTICS_SERVICE_GRPC", "localhost:8004"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to statistics service: %v", err)
	}
	defer statConn.Close()
	statClient := statpb.NewStatisticsServiceClient(statConn)

	r.POST("/api/v1/orders", func(c *gin.Context) {
		var req orderpb.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := orderClient.CreateOrder(context.Background(), &req)
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/orders/:id", func(c *gin.Context) {
		res, err := orderClient.GetOrderByID(context.Background(), &orderpb.GetOrderRequest{
			Id: c.Param("id"),
		})
		handleResponse(c, res, err)
	})

	r.POST("/api/v1/orders/:id/status", func(c *gin.Context) {
		var req orderpb.UpdateOrderStatusRequest
		req.Id = c.Param("id")
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := orderClient.UpdateOrderStatus(context.Background(), &req)
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/orders", func(c *gin.Context) {
		userID := c.Query("user_id")
		page := queryInt(c, "page", 1)
		limit := queryInt(c, "limit", 10)
		res, err := orderClient.ListUserOrders(context.Background(), &orderpb.ListOrdersRequest{
			UserId: userID, Page: int32(page), Limit: int32(limit),
		})
		handleResponse(c, res, err)
	})

	r.POST("/api/v1/products", func(c *gin.Context) {
		var req inventorypb.CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := inventoryClient.CreateProduct(context.Background(), &req)
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/products/:id", func(c *gin.Context) {
		res, err := inventoryClient.GetProductByID(context.Background(), &inventorypb.GetProductRequest{
			Id: c.Param("id"),
		})
		handleResponse(c, res, err)
	})

	r.DELETE("/api/v1/products/:id", func(c *gin.Context) {
		_, err := inventoryClient.DeleteProduct(context.Background(), &inventorypb.DeleteProductRequest{
			Id: c.Param("id"),
		})
		handleResponse(c, gin.H{"message": "deleted"}, err)
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		res, err := inventoryClient.ListProducts(context.Background(), &inventorypb.ListProductsRequest{
			Name:       optional(c.Query("name")),
			CategoryId: optional(c.Query("category_id")),
			Limit:      int32(queryInt(c, "limit", 10)),
			Page:       int32(queryInt(c, "page", 1)),
		})
		handleResponse(c, res, err)
	})

	r.POST("/api/v1/categories", func(c *gin.Context) {
		var req inventorypb.CreateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := inventoryClient.CreateCategory(context.Background(), &req)
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/categories/:id", func(c *gin.Context) {
		res, err := inventoryClient.GetCategoryByID(context.Background(), &inventorypb.GetCategoryRequest{
			Id: c.Param("id"),
		})
		handleResponse(c, res, err)
	})

	r.DELETE("/api/v1/categories/:id", func(c *gin.Context) {
		_, err := inventoryClient.DeleteCategory(context.Background(), &inventorypb.DeleteCategoryRequest{
			Id: c.Param("id"),
		})
		handleResponse(c, gin.H{"message": "deleted"}, err)
	})

	r.GET("/api/v1/categories", func(c *gin.Context) {
		res, err := inventoryClient.ListCategories(context.Background(), &inventorypb.ListCategoriesRequest{
			Name: optional(c.Query("name")),
		})
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/statistics/user-orders/:user_id", func(c *gin.Context) {
		res, err := statClient.GetUserOrdersStatistics(context.Background(), &statpb.UserOrderStatisticsRequest{
			UserId: c.Param("user_id"),
		})
		handleResponse(c, res, err)
	})

	r.GET("/api/v1/statistics/user/:user_id", func(c *gin.Context) {
		res, err := statClient.GetUserStatistics(context.Background(), &statpb.UserStatisticsRequest{
			UserId: c.Param("user_id"),
		})
		handleResponse(c, res, err)
	})

	server := &http.Server{
		Addr:    "0.0.0.0:" + getEnv("HTTP_PORT", "8003"),
		Handler: r,
	}

	go func() {
		log.Println("API Gateway running on", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	waitForShutdown(server)
}

func handleResponse(c *gin.Context, res any, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}
	log.Println("Server exited gracefully")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func queryInt(c *gin.Context, key string, defaultVal int) int {
	val, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return defaultVal
	}
	return val
}

func optional(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}
