package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/order-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, dto dto.OrderCreateDTO) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error)
	GetOrders(ctx context.Context, filter dto.OrderFilterDTO) ([]domain.Order, error)
}

type orderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(repo repository.OrderRepository) *orderUseCase {
	return &orderUseCase{
		orderRepo: repo,
	}
}

func (uc *orderUseCase) CreateOrder(ctx context.Context, dto dto.OrderCreateDTO) (*domain.Order, error) {
	userID, err := primitive.ObjectIDFromHex(dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	productIDs := make([]string, len(dto.Items))
	for i, item := range dto.Items {
		productIDs[i] = item.ProductID
	}

	items := make([]domain.OrderItem, len(dto.Items))
	for i, item := range dto.Items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}
		items[i] = domain.OrderItem{
			ProductID: productID,
			Quantity:  item.Quantity,
		}
	}

	order := &domain.Order{
		ID:     primitive.ObjectID(primitive.NewObjectID()),
		UserID: userID,
		Items:  items,
		Status: domain.OrderStatusPending,
	}

	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (uc *orderUseCase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	orderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (uc *orderUseCase) UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error) {
	orderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	orderStatus := domain.OrderStatus(status)
	switch orderStatus {
	case domain.OrderStatusPending, domain.OrderStatusCompleted, domain.OrderStatusCancelled:
	default:
		return nil, errors.New("invalid order status")
	}

	if err := uc.orderRepo.UpdateOrderStatus(ctx, orderID, orderStatus); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return uc.orderRepo.GetOrderByID(ctx, orderID)
}

func (uc *orderUseCase) GetOrders(ctx context.Context, filter dto.OrderFilterDTO) ([]domain.Order, error) {
	if authUserID, exists := ctx.Value("userID").(string); exists && authUserID != "" {
		if filter.UserID != "" && filter.UserID != authUserID {
			return nil, errors.New("not authorized to view other users' orders")
		}

		if filter.UserID == "" {
			filter.UserID = authUserID
		}
	}

	return uc.orderRepo.GetOrders(ctx, filter)
}
