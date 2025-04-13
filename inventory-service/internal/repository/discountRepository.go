package repository

import (
	"context"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscountRepository interface {
	CreateDiscount(ctx context.Context, discount *domain.Discount) error
	GetAllProductsWithPromotion(ctx context.Context, id primitive.ObjectID) ([]domain.Product, error)
	DeleteDiscount(ctx context.Context, id primitive.ObjectID) error
}

type discountRepository struct {
	collection        *mongo.Collection
	productCollection *mongo.Collection
}

func NewDiscountRepository(db *mongo.Database) *discountRepository {
	return &discountRepository{
		collection:        db.Collection("discounts"),
		productCollection: db.Collection("products"),
	}
}

func (d *discountRepository) CreateDiscount(ctx context.Context, discount *domain.Discount) error {
	_, err := d.collection.InsertOne(ctx, discount)
	return err
}

func (d *discountRepository) DeleteDiscount(ctx context.Context, id primitive.ObjectID) error {
	_, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (d *discountRepository) GetAllProductsWithPromotion(ctx context.Context, id primitive.ObjectID) ([]domain.Product, error) {
	var discount domain.Discount
	err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&discount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var result []domain.Product
	for _, v := range discount.ApplicableProducts {
		var product domain.Product
		err := d.productCollection.FindOne(ctx, bson.M{"_id": v}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		result = append(result, product)
	}
	return result, nil
}
