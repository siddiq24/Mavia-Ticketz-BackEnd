package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/models"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type MovieHandler struct {
	repo repositories.MovieRepository
}

func NewMovieHandler(repo repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

func (h *MovieHandler) GetUpcoming(c *gin.Context) {
	movies, err := h.repo.GetUpcoming()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetPopular(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	movies, err := h.repo.GetPopular(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch popular movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetByFilter(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	genre := c.DefaultQuery("genre", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	movies, err := h.repo.GetByFilter(name, genre, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"page":   page,
		"limit":  limit,
		"result": movies,
	})
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	movie, err := h.repo.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

// GetAllMovies godoc
// @Summary      Get all movies
// @Description  Retrieve all movies (public access)
// @Tags         Movies
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Movie
// @Router       /movies [get]
func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.repo.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": movies})
}

// UpdateMovie godoc
// @Summary      Update movie
// @Description  Admin only - edit movie info
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Movie ID"
// @Param        request body models.UpdateMovieRequest true "Update Movie"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /movies/{id} [patch]
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req models.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ReleaseDate != nil {
		if _, err := time.Parse("2006-01-02", *req.ReleaseDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release date format, use YYYY-MM-DD"})
			return
		}
	}

	if err := h.repo.UpdateMovie(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated successfully"})
}

// DeleteMovie godoc
// @Summary      Delete movie
// @Description  Admin only - delete a movie
// @Tags         Movies
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Movie ID"
// @Success      200 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	if err := h.repo.DeleteMovie(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
