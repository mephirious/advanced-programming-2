package dto

import "time"

type OrderItemDTO struct {
	ProductID string
	Quantity  int
}

type OrderEventDTO struct {
	ID        string
	UserID    string
	Items     []OrderItemDTO
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	EventType string
}

type InventoryEventDTO struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EventType   string
}
