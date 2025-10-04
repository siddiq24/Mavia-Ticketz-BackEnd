package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRouter := router.Group("/movies")
	adminRouter := router.Group("/admin/movies")

	userRepo := repositories.NewUserRepository(db, rdb)
	jwtMiddleware := middlewares.JWTAuthMiddleware(userRepo)
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)

	adminRouter.Use(middlewares.VerifyToken(rdb), jwtMiddleware)

	movieRepo := repositories.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)
	adminHandler := handlers.NewAdminHandler(movieRepo)

	// =======================
	// Public Routes (User)
	// =======================
	movieRouter.GET("/:id", movieHandler.GetMovieByID)
	movieRouter.GET("/upcoming", movieHandler.GetUpcoming)
	movieRouter.GET("/popular", movieHandler.GetPopular)
	movieRouter.GET("", movieHandler.GetByFilter)
	movieRouter.GET("/genres", movieHandler.GetGenres)
	movieRouter.GET("/:id/schedules", scheduleHandler.GetSchedulesById)

	// // serve gambar poster & backdrop
	// movieRouter.GET("/poster/:filename", movieHandler.GetPoster)
	// movieRouter.GET("/backdrop/:filename", movieHandler.GetBackdrop)

	// =======================
	// Admin Routes
	// =======================
	adminRouter.POST("", adminHandler.CreateMovieWithSchedules) // POST /admin/movies
	adminRouter.GET("/uploads/:filename", adminHandler.ServeUploadedFiles)
	adminRouter.GET("", adminHandler.GetMovies)    // GET /admin/movies
	adminRouter.GET("/:id", adminHandler.GetMovie) // GET /admin/movies/:id
	// adminRouter.PATCH("/:id", movieHandler.UpdateMovie)         // PATCH /admin/movies/:id
	adminRouter.DELETE("/:id", adminHandler.DeleteMovie) // DELETE /admin/movies/:id
	// adminRouter.PUT("/:id/backdrop", movieHandler.UpdateBackdrop)
}
