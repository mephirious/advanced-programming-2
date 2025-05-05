package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventID   string             `bson:"event_id"`
	Operation string             `bson:"operation"`
	OrderID   string             `bson:"order_id"`
	UserID    string             `bson:"user_id"`
	Items     []OrderItem        `bson:"items"`
	Total     float64            `bson:"total"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id"`
	Quantity  int     `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type InventoryEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventID   string             `bson:"event_id"`
	Operation string             `bson:"operation"`
	ProductID string             `bson:"product_id"`
	Stock     int32              `bson:"stock"`
	Price     float64            `bson:"price"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserStatistics struct {
	UserID              string      `bson:"user_id"`
	TotalOrders         int         `bson:"total_orders"`
	CompletedOrders     int         `bson:"completed_orders"`
	TotalItemsPurchased int         `bson:"total_items_purchased"`
	AverageOrderValue   float64     `bson:"average_order_value"`
	MostPurchasedItem   string      `bson:"most_purchased_item"`
	HourlyDistribution  map[int]int `bson:"hourly_distribution"`
}
