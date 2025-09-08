package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type ScheduleHandler struct {
	repo repositories.ScheduleRepository
}

func NewScheduleHandler(repo repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{repo: repo}
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	schedules, err := h.repo.GetSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}
