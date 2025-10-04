package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	profileRouter := router.Group("/profile")
	profileRouter.Use(middlewares.VerifyToken(rdb))

	profileRouter.GET("", profileHandler.GetProfile)
	profileRouter.PATCH("", profileHandler.UpdateProfile)
}
