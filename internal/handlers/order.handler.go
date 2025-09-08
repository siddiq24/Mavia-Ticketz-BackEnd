package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/models"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type OrderHandler struct {
	repo repositories.OrderRepository
}

func NewOrderHandler(repo repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  User create order with schedule & seats
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateOrderRequest true "Order Request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ambil user_id dari context (diset oleh middleware VerifyToken)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID int
	switch v := userIDVal.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
		return
	}

	order, err := h.repo.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": order})
}

// GetHistory godoc
// @Summary      Get order history
// @Description  Get all past orders of logged-in user
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /profile/history [get]
func (h *OrderHandler) GetHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int)

	history, err := h.repo.GetOrderHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
