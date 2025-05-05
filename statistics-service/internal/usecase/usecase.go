package usecase

import (
	"context"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/repository"
)

type StatsUseCase interface {
	HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error
	HandleInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error
	GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserStatistics, error)
	GetUserHourlyStatistics(ctx context.Context, userID string) (map[int]int, error)
}

type statsUseCase struct {
	repo repository.StatsRepository
}

func NewStatsUseCase(repo repository.StatsRepository) StatsUseCase {
	return &statsUseCase{repo: repo}
}

func (uc *statsUseCase) HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	return uc.repo.SaveOrderEvent(ctx, event)
}

func (uc *statsUseCase) HandleInventoryEvent(ctx context.Context, event *domain.InventoryEvent) error {
	return uc.repo.SaveInventoryEvent(ctx, event)
}

func (uc *statsUseCase) GetUserOrderStatistics(ctx context.Context, userID string) (*domain.UserStatistics, error) {
	return uc.repo.GetUserOrderStatistics(ctx, userID)
}

func (uc *statsUseCase) GetUserHourlyStatistics(ctx context.Context, userID string) (map[int]int, error) {
	events, err := uc.repo.GetUserOrderStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	hourlyStats := make(map[int]int)
	for _, event := range events {
		if event.Status == "completed" {
			hour := event.CreatedAt.Hour()
			hourlyStats[hour]++
		}
	}
	return hourlyStats, nil
}
