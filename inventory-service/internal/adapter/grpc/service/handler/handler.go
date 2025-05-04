package handler

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
	inventory "github.com/mephirious/advanced-programming-2/inventory-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InventoryHandler struct {
	productUC  usecase.ProductUseCase
	categoryUC usecase.CategoryUseCase
	inventory.UnimplementedInventoryServiceServer
}

func NewInventoryHandler(productUC usecase.ProductUseCase, categoryUC usecase.CategoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		productUC:  productUC,
		categoryUC: categoryUC,
	}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *inventory.CreateProductRequest) (*inventory.Product, error) {
	dto := dto.ProductCreateDTO{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CategoryID:  req.GetCategoryId(),
		Price:       req.GetPrice(),
		Stock:       req.GetStock(),
	}

	product, err := h.productUC.CreateProduct(ctx, dto)
	if err != nil {
		return nil, err
	}

	return mapProductToProto(product), nil
}

func (h *InventoryHandler) GetProductByID(ctx context.Context, req *inventory.GetProductRequest) (*inventory.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	product, err := h.productUC.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapProductToProto(product), nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *inventory.UpdateProductRequest) (*inventory.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	dto := dto.ProductUpdateDTO{
		Name:        optionalString(req.GetName()),
		Description: optionalString(req.GetDescription()),
		CategoryID:  optionalString(req.GetCategoryId()),
		Price:       optionalFloat64(req.GetPrice()),
		Stock:       optionalInt32(req.GetStock()),
	}

	product, err := h.productUC.UpdateProduct(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	return mapProductToProto(product), nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, req *inventory.DeleteProductRequest) (*empty.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := h.productUC.DeleteProduct(ctx, id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, req *inventory.ListProductsRequest) (*inventory.ListProductsResponse, error) {
	dto := dto.ProductFilterDTO{
		Name:       optionalString(req.GetName()),
		CategoryID: optionalString(req.GetCategoryId()),
		MinPrice:   optionalFloat64(req.GetMinPrice()),
		MaxPrice:   optionalFloat64(req.GetMaxPrice()),
		Limit:      req.GetLimit(),
		Page:       req.GetPage(),
		SortBy:     req.GetSortBy(),
		SortOrder:  req.GetSortOrder(),
	}

	products, err := h.productUC.GetAllProducts(ctx, dto)
	if err != nil {
		return nil, err
	}

	var protoProducts []*inventory.Product
	for _, product := range products {
		protoProducts = append(protoProducts, mapProductToProto(&product))
	}

	return &inventory.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// Utility Functions
func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalFloat64(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

func optionalInt32(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

// Proto Mappers
func mapProductToProto(p *domain.Product) *inventory.Product {
	return &inventory.Product{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		CategoryId:  p.CategoryID.Hex(),
		Price:       p.Price,
		Stock:       int32(p.Stock),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

func (h *InventoryHandler) CreateCategory(ctx context.Context, req *inventory.CreateCategoryRequest) (*inventory.Category, error) {
	dto := dto.CategoryCreateDTO{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	category, err := h.categoryUC.CreateCategory(ctx, dto)
	if err != nil {
		return nil, err
	}

	return mapCategoryToProto(category), nil
}

func (h *InventoryHandler) GetCategoryByID(ctx context.Context, req *inventory.GetCategoryRequest) (*inventory.Category, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	category, err := h.categoryUC.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapCategoryToProto(category), nil
}

func (h *InventoryHandler) UpdateCategory(ctx context.Context, req *inventory.UpdateCategoryRequest) (*inventory.Category, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	dto := dto.CategoryUpdateDTO{
		Name:        optionalString(req.GetName()),
		Description: optionalString(req.GetDescription()),
	}

	category, err := h.categoryUC.UpdateCategory(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	return mapCategoryToProto(category), nil
}

func (h *InventoryHandler) DeleteCategory(ctx context.Context, req *inventory.DeleteCategoryRequest) (*empty.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := h.categoryUC.DeleteCategory(ctx, id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h *InventoryHandler) ListCategories(ctx context.Context, req *inventory.ListCategoriesRequest) (*inventory.ListCategoriesResponse, error) {
	categories, err := h.categoryUC.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	var protoCategories []*inventory.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, mapCategoryToProto(&category))
	}

	return &inventory.ListCategoriesResponse{
		Categories: protoCategories,
	}, nil
}

func mapCategoryToProto(c *domain.Category) *inventory.Category {
	return &inventory.Category{
		Id:          c.ID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}
