package domain

import "time"

const (
	OrderStatusCreated   = "created"
	OrderStatusUpdated   = "updated"
	OrderStatusDeleted   = "deleted"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
	OrderStatusPending   = "pending"
)

type EventType int

const (
	CREATED EventType = iota
	UPDATED
	DELETED
)

type OrderEvent struct {
	ID        string      `bson:"id"`
	UserID    string      `bson:"userid"`
	Items     []OrderItem `bson:"items"`
	Total     float64     `bson:"total"`
	Status    string      `bson:"status"`
	CreatedAt time.Time   `bson:"createdat"`
	UpdatedAt time.Time   `bson:"updatedat"`
	EventType string      `bson:"eventtype"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id"`
	Quantity  int     `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type InventoryEvent struct {
	ID          string    `bson:"id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CategoryID  string    `bson:"category_id"`
	Price       float64   `bson:"price"`
	Quantity    int       `bson:"quantity"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	EventType   string    `bson:"event_type"`
}

type UserOrderStatistics struct {
	UserID               string
	TotalOrders          int
	TotalCompletedOrders int
	TotalCancelledOrders int
	OrdersPerHour        map[int]int
}

type UserStatistics struct {
	UserID         string
	TotalUsers     int
	UserOrderCount int
	MostActiveHour int
}
