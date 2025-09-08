package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type SeatHandler struct {
	repo repositories.SeatRepository
}

func NewSeatHandler(repo repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{repo}
}

func (h *SeatHandler) GetAvailableSeats(c *gin.Context) {
	scheduleIDStr := c.Param("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	seats, err := h.repo.GetAvailableSeats(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch available seats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seats})
}
