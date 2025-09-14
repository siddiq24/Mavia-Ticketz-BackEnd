package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
	"github.com/siddiq24/Tickitz-DB/internal/utils"
)

type AuthHandler struct {
	repo repositories.UserRepository
}

func NewAuthHandler(repo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Profile(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"message":  "authenticated",
		"username": username,
	})
}

// Register godoc
// @Summary      Register user
// @Description  Register a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Register Request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hashed, _ := utils.HashPassword(body.Password)
	user, err := h.repo.CreateUser(body.Username, hashed, "user", body.Email) // default role = user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered successfully", "user": user})
}

// Login godoc
// @Summary      Login user
// @Description  Login with email & password to get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login Request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Tentukan identifier yang akan digunakan
	var identifier string
	if body.Username != "" {
		identifier = body.Username
	} else if body.Email != "" {
		identifier = body.Email
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email is required"})
		return
	}

	user, err := h.repo.GetUserByUsernameOrEmail(identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Logout handler
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// blacklist via repo
	err := h.repo.BlacklistToken(context.Background(), tokenString, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to blacklist token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
