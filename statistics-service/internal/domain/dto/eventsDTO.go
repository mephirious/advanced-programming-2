package dto

import (
	"time"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
)

type OrderEventDTO struct {
	EventID   string         `json:"event_id"`
	Operation string         `json:"operation"`
	OrderID   string         `json:"order_id"`
	UserID    string         `json:"user_id"`
	Items     []OrderItemDTO `json:"items"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
}

type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func ToDomainOrderEvent(dto OrderEventDTO) domain.OrderEvent {
	return domain.OrderEvent{
		EventID:   dto.EventID,
		Operation: dto.Operation,
		OrderID:   dto.OrderID,
		UserID:    dto.UserID,
		Items:     ToDomainOrderItems(dto.Items),
		Total:     dto.Total,
		Status:    dto.Status,
		CreatedAt: dto.CreatedAt,
	}
}

type InventoryEventDTO struct {
	EventID   string    `json:"event_id"`
	Operation string    `json:"operation"`
	ProductID string    `json:"product_id"`
	Stock     int32     `json:"stock"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToDomainInventoryEvent(dto InventoryEventDTO) domain.InventoryEvent {
	return domain.InventoryEvent{
		EventID:   dto.EventID,
		Operation: dto.Operation,
		ProductID: dto.ProductID,
		Stock:     dto.Stock,
		Price:     dto.Price,
		UpdatedAt: dto.UpdatedAt,
	}
}

func ToDomainOrderItems(dtoItems []OrderItemDTO) []domain.OrderItem {
	var items []domain.OrderItem
	for _, dtoItem := range dtoItems {
		items = append(items, domain.OrderItem{
			ProductID: dtoItem.ProductID,
			Quantity:  dtoItem.Quantity,
			Price:     dtoItem.Price,
		})
	}
	return items
}

type UserStatisticsDTO struct {
	UserID              string  `json:"user_id"`
	TotalOrders         int     `json:"total_orders"`
	CompletedOrders     int     `json:"completed_orders"`
	TotalItemsPurchased int     `json:"total_items_purchased"`
	AverageOrderValue   float64 `json:"average_order_value"`
	MostPurchasedItem   string  `json:"most_purchased_item"`
}

func ToDomainUserStatistics(dto UserStatisticsDTO) domain.UserStatistics {
	return domain.UserStatistics{
		UserID:              dto.UserID,
		TotalOrders:         dto.TotalOrders,
		CompletedOrders:     dto.CompletedOrders,
		TotalItemsPurchased: dto.TotalItemsPurchased,
		AverageOrderValue:   dto.AverageOrderValue,
		MostPurchasedItem:   dto.MostPurchasedItem,
	}
}
