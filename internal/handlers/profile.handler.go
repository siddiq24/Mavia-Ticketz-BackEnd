package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/models"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type ProfileHandler struct {
	repo repositories.ProfileRepository
}

func NewProfileHandler(repo repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{repo: repo}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile of logged-in user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ProfileResponse
// @Failure 401 {object} map[string]string
// @Router /auth/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int)

	profile, err := h.repo.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Edit profile info (username, avatar, etc.)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.UpdateProfileRequest true "Update Profile"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /profile [patch]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Birthdate != nil {
		if _, err := time.Parse("2006-01-02", *req.Birthdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birthdate format, use YYYY-MM-DD"})
			return
		}
	}

	if err := h.repo.UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
