package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name               string               `json:"name" bson:"name"`
	Description        string               `json:"description" bson:"description"`
	DiscountPercentage float64              `json:"discount_percentage" bson:"discount_percentage"`
	ApplicableProducts []primitive.ObjectID `json:"applicable_products" bson:"applicable_products"`
	StartDate          time.Time            `json:"start_date" bson:"start_date"`
	EndDate            time.Time            `json:"end_date" bson:"end_date"`
	IsActive           bool                 `json:"is_active" bson:"is_active"`
}
