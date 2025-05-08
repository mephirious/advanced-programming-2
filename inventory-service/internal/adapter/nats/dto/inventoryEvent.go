package dto

import (
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	pb "github.com/mephirious/advanced-programming-2/inventory-service/proto/events"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToInventoryEvent(product *domain.Product, eventType pb.InventoryEventType) *pb.InventoryEvent {
	return &pb.InventoryEvent{
		Id:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		CategoryId:  product.CategoryID.Hex(),
		Price:       product.Price,
		Stock:       int32(product.Stock),
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
		EventType:   eventType,
	}
}
