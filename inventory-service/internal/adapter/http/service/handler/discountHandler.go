package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountHandler struct {
	discountUC usecase.DiscountUseCase
}

func NewDiscountHandler(discountUC usecase.DiscountUseCase) *DiscountHandler {
	return &DiscountHandler{
		discountUC: discountUC,
	}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var req dto.DiscountResponseDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DiscountPercentage < 0 || req.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid number of percentage"})
		return
	}

	discount, err := h.discountUC.CreateDiscount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.MapDiscountToResponseDTO(*discount))
}

func (h *DiscountHandler) GetAllProductsWithPromotion(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	discount, err := h.discountUC.GetAllProductsWithPromotion(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discount)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.discountUC.DeleteDiscount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
