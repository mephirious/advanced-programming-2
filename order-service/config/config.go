package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mephirious/advanced-programming-2/order-service/pkg/mongo"
)

type (
	Config struct {
		Mongo  mongo.Config
		Server Server
	}

	Server struct {
		HTTPServer     HTTPServer
		GatewayService string `env:"GATEWAY_SERVICE" envDefault:"0.0.0.0:8003"`
	}

	HTTPServer struct {
		Port int    `env:"HTTP_PORT,required"`
		Mode string `env:"GIN_MODE" envDefault:"release"`
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	var cfg Config

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT value: %w", err)
	}
	cfg.Server.HTTPServer.Port = portInt

	cfg.Server.HTTPServer.Mode = os.Getenv("GIN_MODE")
	if cfg.Server.HTTPServer.Mode == "" {
		cfg.Server.HTTPServer.Mode = "release"
	}

	cfg.Mongo.Database = os.Getenv("MONGO_DB")
	cfg.Mongo.URI = os.Getenv("MONGO_DB_URI")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")

	cfg.Server.GatewayService = os.Getenv("GATEWAY_SERVICE")
	if cfg.Server.GatewayService == "" {
		cfg.Server.GatewayService = "0.0.0.0:8003"
	}

	return &cfg, nil
}
