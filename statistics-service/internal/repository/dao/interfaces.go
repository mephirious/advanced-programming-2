package dao

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderDAO interface {
	Insert(ctx context.Context, order *domain.OrderEvent) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter bson.M) ([]domain.OrderEvent, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
}

type InventoryDAO interface {
	Insert(ctx context.Context, inventory *domain.InventoryEvent) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter bson.M) ([]domain.InventoryEvent, error)
}
