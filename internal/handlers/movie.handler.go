package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

type MovieHandler struct {
	repo repositories.MovieRepository
}

func NewMovieHandler(repo repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

// func (h *MovieHandler) GetUpcoming(c *gin.Context) {
// 	movies, err := h.repo.GetUpcoming()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming movies"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, movies)
// }

// func (h *MovieHandler) GetPopular(c *gin.Context) {
// 	limitStr := c.DefaultQuery("limit", "10")
// 	limit, _ := strconv.Atoi(limitStr)

// 	movies, err := h.repo.GetPopular(limit)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch popular movies"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, movies)
// }

func (h *MovieHandler) GetByFilter(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	genresQuery := c.QueryArray("genres")
	pageStr := c.DefaultQuery("page", "1")

	page, _ := strconv.Atoi(pageStr)
	var genreIDs []int
	if len(genresQuery) == 1 && strings.Contains(genresQuery[0], ",") {
		// Kalau formatnya "1,2,3"
		genresQuery = strings.Split(genresQuery[0], ",")
	}
	for _, g := range genresQuery {
		if g == "" {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSpace(g))
		if err == nil {
			genreIDs = append(genreIDs, id)
		}
	}

	movies, tot, err := h.repo.GetByFilter(name, genreIDs, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"page":       page,
		"total_page": tot,
		"result":     movies,
	})
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	movie, err := h.repo.GetMovieByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
}
func (h *MovieHandler) GetUpcoming(c *gin.Context) {
	movie, err := h.repo.GetUpcomingMovie(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
}
func (h *MovieHandler) GetPopular(c *gin.Context) {
	movie, err := h.repo.GetPopularMovie(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

func (h *MovieHandler) GetGenres(c *gin.Context) {
	genres, err := h.repo.GetGenres(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": genres,
	})
}

// // GetAllMovies godoc
// // @Summary      Get all movies
// // @Description  Retrieve all movies (public access)
// // @Tags         Movies
// // @Produce      json
// // @Security     BearerAuth
// // @Success      200 {array} models.Movie
// // @Router       /movies [get]
// func (h *MovieHandler) GetAllMovies(c *gin.Context) {
// 	movies, err := h.repo.GetAllMovies()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": movies})
// }

// // UpdateMovie godoc
// // @Summary      Update movie
// // @Description  Admin only - edit movie info
// // @Tags         Movies
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Param        id path int true "Movie ID"
// // @Param        request body models.UpdateMovieRequest true "Update Movie"
// // @Success      200 {object} map[string]interface{}
// // @Failure      400 {object} map[string]interface{}
// // @Failure      403 {object} map[string]interface{}
// // @Router       /movies/{id} [patch]
// func (h *MovieHandler) UpdateMovie(c *gin.Context) {
// 	role, _ := c.Get("role")
// 	if role != "admin" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
// 		return
// 	}

// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
// 		return
// 	}

// 	var req models.UpdateMovieRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Printf("BindJSON error: %+v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	log.Printf("Request body: %+v\n", req)

// 	if req.ReleaseDate != nil {
// 		if _, err := time.Parse("2006-01-02", *req.ReleaseDate); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release date format, use YYYY-MM-DD"})
// 			return
// 		}
// 	}

// 	if err := h.repo.UpdateMovie(id, req); err != nil {
// 		if err.Error() == "movie not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
// 			return
// 		} else {
// 			log.Printf("UpdateMovie error: %+v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "movie updated successfully"})
// }

// // DeleteMovie godoc
// // @Summary      Delete movie
// // @Description  Admin only - delete a movie
// // @Tags         Movies
// // @Produce      json
// // @Security     BearerAuth
// // @Param        id path int true "Movie ID"
// // @Success      200 {object} map[string]interface{}
// // @Failure      403 {object} map[string]interface{}
// // @Router       /movies/{id} [delete]
// func (h *MovieHandler) DeleteMovie(c *gin.Context) {
// 	role, exists := c.Get("role")
// 	if !exists || role != "admin" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
// 		c.Abort()
// 		return
// 	}

// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
// 		c.Abort()
// 		return
// 	}

// 	if err := h.repo.DeleteMovie(id); err != nil {
// 		if err.Error() == "movie not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete movie"})
// 			log.Println(err.Error())
// 		}
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
// }

// func saveFile(file *multipart.FileHeader, folder string) (string, error) {
// 	if file == nil {
// 		return "", nil
// 	}
// 	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
// 	path := filepath.Join("uploads", folder, filename)

// 	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
// 		return "", err
// 	}

// 	if err := saveUploadedFile(file, path); err != nil {
// 		return "", err
// 	}
// 	return path, nil
// }

// // saveUploadedFile mirip gin.Context.SaveUploadedFile tapi standalone
// func saveUploadedFile(file *multipart.FileHeader, dst string) error {
// 	src, err := file.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()

// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, src)
// 	return err
// }

// // CreateMovie godoc
// // @Summary Create a new movie
// // @Description Create a new movie with poster, backdrop upload, and relations to casters & genres
// // @Tags Movies
// // @Accept multipart/form-data
// // @Produce json
// // @Param title formData string true "Movie title"
// // @Param description formData string true "Movie description"
// // @Param release_date formData string true "Release date (YYYY-MM-DD)"
// // @Param duration formData int true "Duration in minutes"
// // @Param directors_id formData int true "Director ID"
// // @Param poster formData file false "Poster image file"
// // @Param backdrop formData file false "Backdrop image file"
// // @Param caster_ids formData []int false "List of caster IDs (comma separated)"
// // @Param genre_ids formData []int false "List of genre IDs (comma separated)"
// // @Success 200 {object} models.Movies
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /movies [post]
// func (h *MovieHandler) CreateMovie(c *gin.Context) {
// 	var body models.MovieUploadBody
// 	if err := c.Bind(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	posterPath, err := saveFile(body.PosterImg, "posters")
// 	if err != nil {
// 		log.Println("poster upload error:", err)
// 	}

// 	backdropPath, err := saveFile(body.BackdropImg, "backdrops")
// 	if err != nil {
// 		log.Println("backdrop upload error:", err)
// 	}

// 	releaseDate, _ := time.Parse("2006-01-02", body.ReleaseDate)

// 	// ambil caster_ids dari form
// 	casterIDs := []int{}
// 	if ids, exists := c.GetPostFormArray("caster_ids"); exists {
// 		for _, id := range ids {
// 			if v, err := strconv.Atoi(id); err == nil {
// 				casterIDs = append(casterIDs, v)
// 			}
// 		}
// 	}

// 	// ambil genre_ids dari form
// 	genreIDs := []int{}
// 	if ids, exists := c.GetPostFormArray("genre_ids"); exists {
// 		for _, id := range ids {
// 			if v, err := strconv.Atoi(id); err == nil {
// 				genreIDs = append(genreIDs, v)
// 			}
// 		}
// 	}

// 	movie := models.Movies{
// 		Title:       body.Title,
// 		Description: body.Description,
// 		ReleaseDate: releaseDate,
// 		Duration:    borepositories
// 		PosterImg:   posterPath,
// 		BackdropImg: backdropPath,
// 		DirectorsId: body.DirectorsId,
// 		Rating:      10,
// 	}

// 	ctx := context.Background()
// 	created, err := h.repo.CreateMovie(ctx, movie, casterIDs, genreIDs)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create movie"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, created)
// }

// func (h *MovieHandler) UpdateBackdrop(c *gin.Context) {
// 	id := c.Param("id")

// 	file, _ := c.FormFile("backdrop")
// 	backdropPath, err := saveFile(file, "backdrops")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload backdrop"})
// 		return
// 	}

// 	ctx := context.Background()
// 	idStr, _ := strconv.Atoi(id)
// 	updated, err := h.repo.UpdateBackdrop(ctx, backdropPath, idStr)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update backdrop"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, updated)
// }

// func (h *MovieHandler) GetPoster(c *gin.Context) {
// 	filename := c.Param("filename")
// 	filePath := fmt.Sprintf("./uploads/posters/%s", filename)

// 	c.File(filePath)
// }

// func (h *MovieHandler) GetBackdrop(c *gin.Context) {
// 	filename := c.Param("filename")
// 	filePath := fmt.Sprintf("./uploads/backdrops/%s", filename)

// 	c.File(filePath)
// }
