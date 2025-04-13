package usecase

import (
	"context"
	"fmt"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountUseCase interface {
	CreateDiscount(ctx context.Context, dto dto.DiscountResponseDTO) (*domain.Discount, error)
	GetAllProductsWithPromotion(ctx context.Context, id primitive.ObjectID) ([]domain.Product, error)
	DeleteDiscount(ctx context.Context, id primitive.ObjectID) error
}

type discountUseCase struct {
	discountRepo repository.DiscountRepository
}

func NewDiscountUseCase(repo repository.DiscountRepository) *discountUseCase {
	return &discountUseCase{
		discountRepo: repo,
	}
}

func (uc *discountUseCase) CreateDiscount(ctx context.Context, dto dto.DiscountResponseDTO) (*domain.Discount, error) {
	var productIDs []primitive.ObjectID
	for _, v := range dto.ApplicableProducts {
		id, err := primitive.ObjectIDFromHex(v)
		productIDs = append(productIDs, id)
		if err != nil {
			return nil, fmt.Errorf("invalid productIDs in applicable_products")
		}
	}

	discount := &domain.Discount{
		ID:                 primitive.ObjectID(primitive.NewObjectID()),
		Name:               dto.Name,
		Description:        dto.Description,
		DiscountPercentage: dto.DiscountPercentage,
		ApplicableProducts: productIDs,
		StartDate:          dto.StartDate,
		EndDate:            dto.EndDate,
		IsActive:           dto.IsActive,
	}

	if err := uc.discountRepo.CreateDiscount(ctx, discount); err != nil {
		return nil, err
	}

	return discount, nil
}

func (uc *discountUseCase) GetAllProductsWithPromotion(ctx context.Context, id primitive.ObjectID) ([]domain.Product, error) {
	products, err := uc.discountRepo.GetAllProductsWithPromotion(ctx, id)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *discountUseCase) DeleteDiscount(ctx context.Context, id primitive.ObjectID) error {
	return uc.discountRepo.DeleteDiscount(ctx, id)
}
