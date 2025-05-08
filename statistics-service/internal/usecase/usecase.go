package usecase

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/repository"
)

type StatsUseCase interface {
	HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error
	HandleInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error
	GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserOrderStatistics, error)
	GetUserHourlyOrderStatistics(ctx context.Context, userID string) (map[int]int, error)
}

type statsUseCase struct {
	repo repository.StatsRepository
}

func NewStatsUseCase(repo repository.StatsRepository) StatsUseCase {
	return &statsUseCase{repo: repo}
}

func (uc *statsUseCase) HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	if event == nil {
		return nil
	}
	return uc.repo.SaveOrderEvent(ctx, event)
}

func (uc *statsUseCase) HandleInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error {
	if event == nil {
		return nil
	}
	return uc.repo.SaveInventoryEvent(ctx, event)
}

func (uc *statsUseCase) GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserOrderStatistics, error) {
	return uc.repo.GetUserOrderStatistics(ctx, userID)
}

func (uc *statsUseCase) GetUserHourlyOrderStatistics(ctx context.Context, userID string) (map[int]int, error) {
	orders, err := uc.repo.GetUserOrderEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	hourlyStats := make(map[int]int)
	for _, order := range orders {
		if order.Status == domain.OrderStatusCompleted {
			hour := order.CreatedAt.Hour()
			hourlyStats[hour]++
		}
	}

	return hourlyStats, nil
}
