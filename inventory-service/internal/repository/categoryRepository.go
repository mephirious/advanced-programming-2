package repository

import (
	"context"
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id primitive.ObjectID) error
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) *categoryRepository {
	return &categoryRepository{
		collection: db.Collection("categories"),
	}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, category)
	return err
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	category.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": category.ID},
		bson.M{"$set": category},
	)
	return err
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
