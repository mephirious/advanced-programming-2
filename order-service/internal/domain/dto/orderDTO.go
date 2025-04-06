package dto

import (
	"time"

	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
)

type OrderCreateDTO struct {
	UserID string         `json:"user_id" binding:"required"`
	Items  []OrderItemDTO `json:"items" binding:"required,min=1"`
}

type OrderItemDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type OrderUpdateDTO struct {
	Status *string `json:"status" binding:"omitempty,oneof=pending completed cancelled"`
}

type OrderFilterDTO struct {
	UserID string `form:"user_id"`
	Limit  int    `form:"limit,default=20"`
	Page   int    `form:"page,default=1"`
	Status string `form:"status"`
}

type OrderResponseDTO struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Items     []OrderItemRespDTO `json:"items"`
	Total     float64            `json:"total"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type OrderItemRespDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func MapOrderToResponseDTO(o domain.Order) OrderResponseDTO {
	items := make([]OrderItemRespDTO, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemRespDTO{
			ProductID: item.ProductID.Hex(),
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return OrderResponseDTO{
		ID:        o.ID.Hex(),
		UserID:    o.UserID.Hex(),
		Items:     items,
		Total:     o.Total,
		Status:    string(o.Status),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
