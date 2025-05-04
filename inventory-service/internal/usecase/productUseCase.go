package usecase

import (
	"context"
	"fmt"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, dto dto.ProductCreateDTO) (*domain.Product, error)
	GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id primitive.ObjectID, dto dto.ProductUpdateDTO) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	GetAllProducts(ctx context.Context, filter dto.ProductFilterDTO) ([]domain.Product, error)
}

type productUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *productUseCase {
	return &productUseCase{
		productRepo: repo,
	}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, dto dto.ProductCreateDTO) (*domain.Product, error) {
	categoryObjectID, err := primitive.ObjectIDFromHex(dto.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_id: %w", err)
	}

	product := &domain.Product{
		ID:          primitive.ObjectID(primitive.NewObjectID()),
		Name:        dto.Name,
		Description: dto.Description,
		CategoryID:  categoryObjectID,
		Price:       dto.Price,
		Stock:       dto.Stock,
	}

	if err := uc.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *productUseCase) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	product, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id primitive.ObjectID, dto dto.ProductUpdateDTO) (*domain.Product, error) {
	product, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		product.Name = *dto.Name
	}
	if dto.Description != nil {
		product.Description = *dto.Description
	}
	if dto.CategoryID != nil {
		categoryID, err := primitive.ObjectIDFromHex(*dto.CategoryID)
		if err != nil {
			return nil, err
		}
		product.CategoryID = categoryID
	}
	if dto.Price != nil {
		product.Price = *dto.Price
	}
	if dto.Stock != nil {
		product.Stock = *dto.Stock
	}

	err = uc.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	p1roduct, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return p1roduct, nil
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return uc.productRepo.DeleteProduct(ctx, id)
}

func (uc *productUseCase) GetAllProducts(ctx context.Context, filter dto.ProductFilterDTO) ([]domain.Product, error) {
	return uc.productRepo.GetAllProducts(ctx, filter)
}
