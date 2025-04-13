package dto

import (
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
)

type DiscountResponseDTO struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ApplicableProducts []string  `json:"applicable_products"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	IsActive           bool      `json:"is_active"`
}

func MapDiscountToResponseDTO(d domain.Discount) DiscountResponseDTO {
	productsID := make([]string, 0)
	for _, v := range d.ApplicableProducts {
		productsID = append(productsID, v.Hex())
	}
	return DiscountResponseDTO{
		ID:                 d.ID.Hex(),
		Name:               d.Name,
		Description:        d.Description,
		DiscountPercentage: d.DiscountPercentage,
		ApplicableProducts: productsID,
		StartDate:          d.StartDate,
		EndDate:            d.EndDate,
		IsActive:           d.IsActive,
	}
}
