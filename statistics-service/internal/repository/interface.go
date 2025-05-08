package repository

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
)

type StatsRepository interface {
	SaveOrderEvent(ctx context.Context, event *domain.OrderEvent) error
	SaveInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error
	GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserOrderStatistics, error)
	GetUserOrderEvents(ctx context.Context, userID string) ([]*domain.OrderEvent, error)
}
