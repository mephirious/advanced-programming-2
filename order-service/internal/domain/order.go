package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Price     float64            `json:"price" bson:"price"`
}

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items     []OrderItem        `json:"items" bson:"items"`
	Total     float64            `json:"total" bson:"total"`
	Status    OrderStatus        `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
