package dto

import (
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
)

type ProductCreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	CategoryID  string  `json:"category_id" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
}

type ProductUpdateDTO struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

type ProductFilterDTO struct {
	Name       *string  `form:"name"`
	CategoryID *string  `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Limit      int      `form:"limit,default=20"`
	Page       int      `form:"page,default=1"`
	SortBy     string   `form:"sort_by"`
	SortOrder  string   `form:"sort_order"`
}

type ProductResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MapProductToResponseDTO(p domain.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		CategoryID:  p.CategoryID.Hex(),
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
