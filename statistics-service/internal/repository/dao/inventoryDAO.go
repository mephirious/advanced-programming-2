package dao

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type inventoryDAO struct {
	collection *mongo.Collection
}

func NewInventoryDAO(db *mongo.Database) InventoryDAO {
	return &inventoryDAO{
		collection: db.Collection("inventory_events"),
	}
}

func (dao *inventoryDAO) Insert(ctx context.Context, inventory *domain.InventoryEvent) (*mongo.InsertOneResult, error) {
	return dao.collection.InsertOne(ctx, inventory)
}

func (dao *inventoryDAO) Find(ctx context.Context, filter bson.M) ([]domain.InventoryEvent, error) {
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var inventories []domain.InventoryEvent
	if err = cursor.All(ctx, &inventories); err != nil {
		return nil, err
	}
	return inventories, nil
}
