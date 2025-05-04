package handler

import (
	"context"
	"fmt"

	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/order-service/internal/usecase"
	orderpb "github.com/mephirious/advanced-programming-2/order-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderUC usecase.OrderUseCase
	orderpb.UnimplementedOrderServiceServer
}

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	items := make([]dto.OrderItemDTO, len(req.GetItems()))
	for i, item := range req.GetItems() {
		items[i] = dto.OrderItemDTO{
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		}
	}

	dto := dto.OrderCreateDTO{
		UserID: req.GetUserId(),
		Items:  items,
	}

	order, err := h.orderUC.CreateOrder(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Order: mapOrderToProto(order),
	}, nil
}

func (h *OrderHandler) GetOrderByID(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	order, err := h.orderUC.GetOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Order: mapOrderToProto(order),
	}, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	order, err := h.orderUC.UpdateOrderStatus(ctx, req.GetId(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Order: mapOrderToProto(order),
	}, nil
}

func (h *OrderHandler) ListUserOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	filter := dto.OrderFilterDTO{
		UserID: req.GetUserId(),
		Page:   req.GetPage(),
		Limit:  req.GetLimit(),
	}

	orders, err := h.orderUC.GetOrders(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	var orderResponses []*orderpb.Order
	for _, order := range orders {
		orderResponses = append(orderResponses, mapOrderToProto(&order))
	}

	return &orderpb.ListOrdersResponse{
		Orders: orderResponses,
	}, nil
}

func mapOrderToProto(o *domain.Order) *orderpb.Order {
	items := make([]*orderpb.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = &orderpb.OrderItem{
			ProductId: item.ProductID.Hex(),
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &orderpb.Order{
		Id:        o.ID.Hex(),
		UserId:    o.UserID.Hex(),
		Items:     items,
		Status:    mapOrderStatusToProto(o.Status),
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}

func mapOrderStatusToProto(status domain.OrderStatus) orderpb.OrderStatus {
	switch status {
	case domain.OrderStatusPending:
		return orderpb.OrderStatus_PENDING
	case domain.OrderStatusCompleted:
		return orderpb.OrderStatus_COMPLETED
	case domain.OrderStatusCancelled:
		return orderpb.OrderStatus_CANCELLED
	default:
		return orderpb.OrderStatus_PENDING
	}
}
