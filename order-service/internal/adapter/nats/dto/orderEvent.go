package dto

import (
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
	pb "github.com/mephirious/advanced-programming-2/order-service/proto/events"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToOrderEvent(order *domain.Order, eventType pb.OrderEventType) *pb.OrderEvent {
	orderItems := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		orderItems[i] = &pb.OrderItem{
			ProductId: item.ProductID.Hex(),
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &pb.OrderEvent{
		Id:        order.ID.Hex(),
		UserId:    order.UserID.Hex(),
		Items:     orderItems,
		Total:     order.Total,
		Status:    domainStatusToProtoStatus(order.Status),
		CreatedAt: timestamppb.New(order.CreatedAt),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
		EventType: eventType,
	}
}

func domainStatusToProtoStatus(status domain.OrderStatus) pb.OrderStatus {
	switch status {
	case domain.OrderStatusPending:
		return pb.OrderStatus_S_PENDING
	case domain.OrderStatusCompleted:
		return pb.OrderStatus_S_COMPLETED
	case domain.OrderStatusCancelled:
		return pb.OrderStatus_S_CANCELLED
	default:
		return pb.OrderStatus_S_PENDING
	}
}
