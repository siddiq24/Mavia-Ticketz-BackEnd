package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	profileRouter := router.Group("/profile")
	profileRouter.Use(middlewares.VerifyToken)

	profileRouter.GET("", profileHandler.GetProfile)
	profileRouter.PATCH("", profileHandler.UpdateProfile)
}
