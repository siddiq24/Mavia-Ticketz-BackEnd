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
	movieRepo := repositories.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movieRouter := router.Group("/movies")
	adminRouter := router.Group("/admin")
	adminRouter.Use(middlewares.VerifyToken)

	{
		movieRouter.GET("/upcoming", movieHandler.GetUpcoming)
		movieRouter.GET("/popular", movieHandler.GetPopular)
		movieRouter.GET("", movieHandler.GetByFilter)
		movieRouter.GET("/:id", movieHandler.GetMovieByID)
		adminRouter.PATCH("/:id", movieHandler.UpdateMovie)
		adminRouter.DELETE("/:id", movieHandler.DeleteMovie)
	}
}
