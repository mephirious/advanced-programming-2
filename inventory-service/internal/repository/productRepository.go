package repository

import (
	"context"
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	GetAllProducts(ctx context.Context, filter dto.ProductFilterDTO) ([]domain.Product, error)
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *productRepository {
	return &productRepository{
		collection: db.Collection("products"),
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepository) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": product.ID},
		bson.M{"$set": product},
	)
	return err
}

func (r *productRepository) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepository) GetAllProducts(ctx context.Context, filter dto.ProductFilterDTO) ([]domain.Product, error) {
	query := bson.M{}
	if filter.Name != nil {
		query["name"] = bson.M{"$regex": *filter.Name, "$options": "i"}
	}
	if filter.CategoryID != nil {
		query["category_id"] = *filter.CategoryID
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceQuery := bson.M{}
		if filter.MinPrice != nil {
			priceQuery["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceQuery["$lte"] = *filter.MaxPrice
		}
		query["price"] = priceQuery
	}

	opts := options.Find()
	opts.SetSkip(int64((filter.Page - 1) * filter.Limit))
	opts.SetLimit(int64(filter.Limit))

	if filter.SortBy != "" {
		sortOrder := 1
		if filter.SortOrder == "desc" {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{Key: filter.SortBy, Value: sortOrder}})
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
