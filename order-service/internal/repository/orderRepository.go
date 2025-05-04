package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id primitive.ObjectID) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error
	GetOrders(ctx context.Context, filter dto.OrderFilterDTO) ([]domain.Order, error)
}

type orderRepository struct {
	collection        *mongo.Collection
	productCollection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *orderRepository {
	return &orderRepository{
		collection:        db.Collection("orders"),
		productCollection: db.Collection("products"),
	}
}
func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	productIDs := make([]primitive.ObjectID, len(order.Items))
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
	}

	cursor, err := r.productCollection.Find(ctx, bson.M{"_id": bson.M{"$in": productIDs}})
	if err != nil {
		return fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	productPrices := make(map[primitive.ObjectID]float64)

	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return fmt.Errorf("failed to decode product: %w", err)
		}
		productPrices[product.ID] = product.Price
	}

	for i, item := range order.Items {
		if price, found := productPrices[item.ProductID]; found {
			order.Items[i].Price = price
		} else {
			return fmt.Errorf("price not found for product ID %v", item.ProductID)
		}
	}

	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	_, err = r.collection.InsertOne(ctx, order)
	return err
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id primitive.ObjectID) (*domain.Order, error) {
	var order domain.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		}},
	)
	return err
}

func (r *orderRepository) GetOrders(ctx context.Context, filter dto.OrderFilterDTO) ([]domain.Order, error) {
	query := bson.M{}

	if filter.UserID != "" {
		userID, err := primitive.ObjectIDFromHex(filter.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
		query["user_id"] = userID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	opts := options.Find()
	opts.SetSkip(int64((filter.Page - 1) * filter.Limit))
	opts.SetLimit(int64(filter.Limit))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
