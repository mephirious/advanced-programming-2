package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler struct {
	categoryUC usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUC usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUC: categoryUC,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryUC.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.MapCategoryToResponseDTO(*category))
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.categoryUC.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MapCategoryToResponseDTO(*category))
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req dto.CategoryUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryUC.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MapCategoryToResponseDTO(*category))
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.categoryUC.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryUC.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.CategoryResponseDTO, len(categories))
	for i, category := range categories {
		response[i] = dto.MapCategoryToResponseDTO(category)
	}

	c.JSON(http.StatusOK, response)
}
