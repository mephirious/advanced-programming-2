package repository

import (
	"context"
	"log"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStatsRepository struct {
	orderCol     *mongo.Collection
	inventoryCol *mongo.Collection
}

func NewMongoStatsRepository(db *mongo.Database) StatsRepository {
	return &mongoStatsRepository{
		orderCol:     db.Collection("order_events"),
		inventoryCol: db.Collection("inventory_events"),
	}
}

func (r *mongoStatsRepository) SaveOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	result, err := r.orderCol.InsertOne(ctx, event)
	if err != nil {
		log.Printf("Failed to save order event: %v", err)
		return err
	}
	log.Printf("Saved order event: %v", result.InsertedID)
	return nil
}

func (r *mongoStatsRepository) SaveInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error {
	result, err := r.inventoryCol.InsertOne(ctx, event)
	if err != nil {
		log.Printf("Failed to save inventory event: %v", err)
		return err
	}
	log.Printf("Saved inventory event: %v", result.InsertedID)
	return nil
}

func (r *mongoStatsRepository) GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserOrderStatistics, error) {
	stats := &domain.UserOrderStatistics{
		UserID:        userID,
		OrdersPerHour: make(map[int]int),
	}

	filter := bson.M{"userid": userID}
	cursor, err := r.orderCol.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to query order events: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var evt domain.OrderEvent
		if err := cursor.Decode(&evt); err != nil {
			log.Printf("Failed to decode order event: %v", err)
			continue
		}
		stats.TotalOrders++
		hour := evt.CreatedAt.Hour()
		stats.OrdersPerHour[hour]++

		switch evt.Status {
		case "S_COMPLETED":
			stats.TotalCompletedOrders++
		case "S_CANCELLED":
			stats.TotalCancelledOrders++
		}
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return stats, nil
}

func (r *mongoStatsRepository) GetUserOrderEvents(ctx context.Context, userID string) ([]*domain.OrderEvent, error) {
	filter := bson.M{"userid": userID}
	cursor, err := r.orderCol.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to query order events: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*domain.OrderEvent
	for cursor.Next(ctx) {
		var evt domain.OrderEvent
		if err := cursor.Decode(&evt); err != nil {
			log.Printf("Failed to decode order event: %v", err)
			continue
		}
		events = append(events, &evt)
	}
	return events, nil
}
