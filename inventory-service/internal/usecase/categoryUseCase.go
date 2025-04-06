package usecase

import (
	"context"
	"errors"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, dto dto.CategoryCreateDTO) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id primitive.ObjectID, dto dto.CategoryUpdateDTO) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id primitive.ObjectID) error
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(repo repository.CategoryRepository) *categoryUseCase {
	return &categoryUseCase{
		categoryRepo: repo,
	}
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, dto dto.CategoryCreateDTO) (*domain.Category, error) {
	existing, err := uc.categoryRepo.GetCategoryByName(ctx, dto.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &domain.Category{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := uc.categoryRepo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *categoryUseCase) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id primitive.ObjectID, dto dto.CategoryUpdateDTO) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	if dto.Name != nil && *dto.Name != category.Name {
		existing, err := uc.categoryRepo.GetCategoryByName(ctx, *dto.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("another category with this name already exists")
		}
	}

	if dto.Name != nil {
		category.Name = *dto.Name
	}
	if dto.Description != nil {
		category.Description = *dto.Description
	}

	if err := uc.categoryRepo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return uc.categoryRepo.DeleteCategory(ctx, id)
}

func (uc *categoryUseCase) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.GetAllCategories(ctx)
}
