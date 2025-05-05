package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Database string
	URI      string
	Username string
	Password string
}

type DB struct {
	Connection *mongo.Database
	Client     *mongo.Client
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connection to mongoDB Error: %w", err)
	}

	db := &DB{
		Connection: client.Database(cfg.Database),
		Client:     client,
	}

	err = db.Client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping connection mongoDB Error: %w", err)
	}

	go reconnectOnFailure(ctx, db.Client, clientOptions)

	return db, nil
}

func reconnectOnFailure(ctx context.Context, client *mongo.Client, opts *options.ClientOptions) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := client.Ping(ctx, nil)
			if err != nil {
				log.Printf("lost connection to mongoDB: %v", err)
				newClient, _ := mongo.Connect(ctx, opts)
				*client = *newClient
			}
		case <-ctx.Done():
			return
		}
	}
}
