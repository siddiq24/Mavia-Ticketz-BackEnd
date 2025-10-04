package handlers

import (
	"net/http"
	"strconv"

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

func (h *ScheduleHandler) GetSchedulesById(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	// call repo
	schedules, err := h.repo.GetSchedulesById(c, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"movie_id":  movieID,
		"schedules": schedules,
	})
}
