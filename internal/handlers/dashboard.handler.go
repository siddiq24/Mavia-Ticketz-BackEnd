package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type DashboardHandler struct {
	repo repositories.DashboardRepository
}

func NewDashboardHandler(repo repositories.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	yearParam := c.Query("year")
	if yearParam == "" {
		yearParam = "2025" // default
	}
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	movieSales, err := h.repo.GetTicketSalesByMovie(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categorySales, err := h.repo.GetTicketSalesByCategoryLocation(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sales_by_movie":             movieSales,
		"sales_by_category_location": categorySales,
	})
}
