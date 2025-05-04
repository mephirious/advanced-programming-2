package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/mongo"
)

type (
	Config struct {
		Mongo  mongo.Config
		Server Server
	}

	Server struct {
		GRPCServer GRPCServer
	}

	GRPCServer struct {
		Port int    `env:"GRPC_PORT,required"`
		Mode string `env:"GIN_MODE" envDefault:"release"`
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	var cfg Config

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT value: %w", err)
	}
	cfg.Server.GRPCServer.Port = portInt

	cfg.Server.GRPCServer.Mode = os.Getenv("GIN_MODE")
	if cfg.Server.GRPCServer.Mode == "" {
		cfg.Server.GRPCServer.Mode = "release"
	}

	cfg.Mongo.Database = os.Getenv("MONGO_DB")
	cfg.Mongo.URI = os.Getenv("MONGO_DB_URI")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")

	return &cfg, nil
}
