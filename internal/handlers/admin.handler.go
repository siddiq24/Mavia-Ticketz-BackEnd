package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/models"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
	"github.com/siddiq24/Tickitz-DB/internal/utils"
)

type AdminHandler struct {
	movieRepo repositories.MovieRepository
}

func NewAdminHandler(movieRepo repositories.MovieRepository) *AdminHandler {
	return &AdminHandler{movieRepo: movieRepo}
}

// Get all Movies
func (h *AdminHandler) GetMovies(c *gin.Context) {
	// movies, err := h.movieRepo.GetAllMovies()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"movies": movies})
}

// Get Movie by ID
func (h *AdminHandler) GetMovieByID(c *gin.Context) {
	// id := c.Param("id")
	// movieID := 0
	// _, err := fmt.Sscan(id, &movieID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	// 	return
	// }

	// movie, err := h.movieRepo.GetMovieByID(movieID)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"movie": movie})
}
func (h *AdminHandler) CreateMovieWithSchedules(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
		c.JSON(http.StatusBadRequest, utils.ResponseError("Failed to parse form", err.Error()))
		return
	}

	// Handle file uploads
	var posterPath string
	var backdropPath string

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError("Failed to create upload directory", err.Error()))
		return
	}

	// Handle poster upload
	poster, err := c.FormFile("poster")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Poster is required", err.Error()))
		return
	}
	posterPath = fmt.Sprintf("uploads/posters/%d_%s", time.Now().UnixNano(), poster.Filename)
	if err := c.SaveUploadedFile(poster, posterPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError("Failed to save poster", err.Error()))
		return
	}
	posterPath = filepath.Base(posterPath)

	// Handle backdrop upload
	backdrop, err := c.FormFile("backdrop")
	if err != nil {
		// Clean up poster if backdrop fails
		os.Remove(filepath.Join("uploads", posterPath))
		c.JSON(http.StatusBadRequest, utils.ResponseError("Backdrop is required", err.Error()))
		return
	}
	backdropPath = fmt.Sprintf("uploads/backdrops/%d_%s", time.Now().UnixNano(), backdrop.Filename)
	if err := c.SaveUploadedFile(backdrop, backdropPath); err != nil {
		// Clean up poster if backdrop fails
		os.Remove(filepath.Join("uploads", posterPath))
		c.JSON(http.StatusInternalServerError, utils.ResponseError("Failed to save backdrop :", err.Error()))
		return
	}
	backdropPath = filepath.Base(backdropPath)

	// Parse form values
	movieReq := models.MovieRequest{
		PosterImg:   posterPath,
		BackdropImg: backdropPath,
	}

	// Parse text fields
	movieReq.Title = c.PostForm("title")
	movieReq.Description = c.PostForm("description")

	// Parse release date
	releaseDateStr := c.PostForm("release_date")
	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)

	if err != nil {
		cleanupFiles(posterPath, backdropPath)
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid release date format", "Use YYYY-MM-DD format"))
		return
	}
	movieReq.ReleaseDate = releaseDate

	// Parse duration
	if durationStr := c.PostForm("duration"); durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			cleanupFiles(posterPath, backdropPath)
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid duration", err.Error()))
			return
		}
		movieReq.Duration = duration
	}

	// Parse director name
	movieReq.DirectorName = c.PostForm("director_name")

	// Parse rating
	if ratingStr := c.PostForm("rating"); ratingStr != "" {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			cleanupFiles(posterPath, backdropPath)
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid rating", err.Error()))
			return
		}
		movieReq.Rating = rating
	}

	// Parse array fields
	movieReq.CasterIDs = parseIntArray(c.PostFormArray("caster_ids[]"))
	movieReq.GenreIDs = parseIntArray(c.PostFormArray("genre_ids[]"))

	// Parse schedules from JSON
	schedulesJSON := c.PostForm("schedules")
	if schedulesJSON != "" {
		var schedules []models.ScheduleInput
		if err := json.Unmarshal([]byte(schedulesJSON), &schedules); err != nil {
			cleanupFiles(posterPath, backdropPath)
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid schedules format", err.Error()))
			return
		}
		movieReq.Schedules = schedules
	}

	// Validate request
	if err := utils.ValidateStruct(movieReq); err != nil {
		cleanupFiles(posterPath, backdropPath)
		c.JSON(http.StatusBadRequest, utils.ResponseError("Validation error", err.Error()))
		return
	}

	movieID, err := h.movieRepo.CreateMovieWithSchedules(c.Request.Context(), movieReq)
	if err != nil {
		cleanupFiles(posterPath, backdropPath)
		c.JSON(http.StatusInternalServerError, utils.ResponseError("Failed to create movie", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseSuccess("Movie created successfully", gin.H{
		"movie_id":     movieID,
		"poster_url":   posterPath,
		"backdrop_url": backdropPath,
	}))
}

// Helper function to cleanup uploaded files
func cleanupFiles(posterPath, backdropPath string) {
	os.Remove(filepath.Join("uploads", posterPath))
	os.Remove(filepath.Join("uploads", backdropPath))
}

// Helper function to parse string array to int array
func parseIntArray(strArr []string) []int {
	var intArr []int
	for _, str := range strArr {
		if num, err := strconv.Atoi(str); err == nil {
			intArr = append(intArr, num)
		}
	}
	return intArr
}

// CreateMovieWithSchedules godoc
// @Summary      Create a new movie with schedules
// @Description  Create a new movie with details, poster & backdrop upload, and schedules.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        title          formData string  true  "Movie Title"
// @Param        description    formData string  true  "Movie Description"
// @Param        release_date   formData string  true  "Release Date (YYYY-MM-DD)"
// @Param        duration       formData int     true  "Duration in minutes"
// @Param        director_id    formData int     true  "Director ID"
// @Param        rating         formData number  false "Movie Rating (0-10)"
// @Param        genre_ids      formData []int   false "Genre IDs (multiple allowed)"
// @Param        caster_ids     formData []int   false "Caster IDs (multiple allowed)"
// @Param        schedules      formData string  false "Schedules JSON string"
// @Param        poster         formData file    false "Poster image"
// @Param        backdrop       formData file    false "Backdrop image"
// @Success      201 {object} map[string]interface{} "Movie created successfully"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      500 {object} map[string]interface{} "Internal Server Error"
// @Router       /admin/movies [post]
func (h *AdminHandler) GetMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid movie ID", err.Error()))
		return
	}

	movie, err := h.movieRepo.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("Movie not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Movie retrieved successfully", movie))
}

func (h *AdminHandler) ServeUploadedFiles(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("uploads", filename)

	// Security check to prevent directory traversal
	if filepath.Base(filePath) != filename {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid filename", "Invalid file path"))
		return
	}

	c.File(filePath)
}

func (h *AdminHandler) DeleteMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	err = h.movieRepo.DeleteMovieById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "movie deleted successfully",
	})
}

func (h *AdminHandler) UpdateMovie(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Failed to parse form", err.Error()))
		return
	}

	// var posterPath string
	// var backdropPath string

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError("failed to create upload directory", err.Error()))
	}
}
