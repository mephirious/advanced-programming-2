package dto

import (
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
)

type CategoryCreateDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CategoryUpdateDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CategoryResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MapCategoryToResponseDTO(c domain.Category) CategoryResponseDTO {
	return CategoryResponseDTO{
		ID:          c.ID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
