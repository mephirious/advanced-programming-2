package dao

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderDAO struct {
	collection *mongo.Collection
}

func NewOrderDAO(db *mongo.Database) OrderDAO {
	return &orderDAO{
		collection: db.Collection("order_events"),
	}
}

func (dao *orderDAO) Insert(ctx context.Context, order *domain.OrderEvent) (*mongo.InsertOneResult, error) {
	return dao.collection.InsertOne(ctx, order)
}

func (dao *orderDAO) Find(ctx context.Context, filter bson.M) ([]domain.OrderEvent, error) {
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var orders []domain.OrderEvent
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (dao *orderDAO) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return dao.collection.Aggregate(ctx, pipeline)
}
