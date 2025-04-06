package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/order-service/internal/usecase"
)

type OrderHandler struct {
	orderUC usecase.OrderUseCase
}

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.OrderCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderUC.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.MapOrderToResponseDTO(*order))
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	order, err := h.orderUC.GetOrderByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MapOrderToResponseDTO(*order))
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req dto.OrderUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderUC.UpdateOrderStatus(c.Request.Context(), c.Param("id"), *req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MapOrderToResponseDTO(*order))
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	var filter dto.OrderFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authUserID, exists := c.Get("userID"); exists {
		filter.UserID = authUserID.(string)
	}

	orders, err := h.orderUC.GetOrders(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.OrderResponseDTO, len(orders))
	for i, order := range orders {
		response[i] = dto.MapOrderToResponseDTO(order)
	}

	c.JSON(http.StatusOK, response)
}
