package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieRepo := repositories.NewMovieRepository(db)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movieRouter := router.Group("/movies")
	movieRouter.Use(middlewares.VerifyToken)

	{
		movieRouter.GET("/upcoming", movieHandler.GetUpcoming)
		movieRouter.GET("/popular", movieHandler.GetPopular)
		movieRouter.GET("", movieHandler.GetByFilter)
		movieRouter.GET("/:id", movieHandler.GetMovieByID)
		movieRouter.PATCH("/:id", movieHandler.UpdateMovie)
		movieRouter.DELETE("/:id", movieHandler.DeleteMovie)

	}
}
