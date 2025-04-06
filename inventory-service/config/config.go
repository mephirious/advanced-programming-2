package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/mongo"
)

type (
	Config struct {
		Mongo  mongo.Config
		Server Server
	}

	Server struct {
		HTTPServer HTTPServer
	}

	HTTPServer struct {
		Port int    `env:"HTTP_PORT,required"`
		Mode string `env:"GIN_MODE" envDefault:"release"`
	}
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}
