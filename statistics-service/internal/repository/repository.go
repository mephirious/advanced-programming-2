package repository

import (
	"context"
	"fmt"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsRepository struct {
	orderDAO     dao.OrderDAO
	inventoryDAO dao.InventoryDAO
}

func NewStatsRepository(orderDAO dao.OrderDAO, inventoryDAO dao.InventoryDAO) *StatsRepository {
	return &StatsRepository{
		orderDAO:     orderDAO,
		inventoryDAO: inventoryDAO,
	}
}

func (r *StatsRepository) SaveOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	_, err := r.orderDAO.Insert(ctx, event)
	return err
}

func (r *StatsRepository) SaveInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error {
	_, err := r.inventoryDAO.Insert(ctx, event)
	return err
}

func (r *StatsRepository) GetUserOrderStats(ctx context.Context, userID string) ([]domain.OrderEvent, error) {
	return r.orderDAO.Find(ctx, bson.M{"user_id": userID})
}

func (r *StatsRepository) GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserStatistics, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"user_id", userID},
			{"status", "completed"},
		}}},
		{{"$unwind", "$items"}},
		{{"$group", bson.D{
			{"_id", "$items.product_id"},
			{"total_quantity", bson.D{{"$sum", "$items.quantity"}}},
			{"total_orders", bson.D{{"$sum", 1}}},
			{"total_spent", bson.D{{"$sum", "$total"}}},
		}}},
		{{"$sort", bson.D{{"total_quantity", -1}}}},
		{{"$limit", 1}},
	}

	cursor, err := r.orderDAO.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		ProductID   string  `bson:"_id"`
		TotalQty    int     `bson:"total_quantity"`
		TotalOrders int     `bson:"total_orders"`
		TotalSpent  float64 `bson:"total_spent"`
	}

	stats := &domain.UserStatistics{
		UserID: userID,
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode aggregation result: %w", err)
		}
		stats.MostPurchasedItem = result.ProductID
		stats.TotalItemsPurchased = result.TotalQty
		if result.TotalOrders > 0 {
			stats.AverageOrderValue = result.TotalSpent / float64(result.TotalOrders)
		}
	}

	totalOrders, err := r.orderDAO.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}
	stats.TotalOrders = len(totalOrders)

	return stats, nil
}
