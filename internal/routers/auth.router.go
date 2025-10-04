package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	repo := repositories.NewUserRepository(db, rdb)
	handler := handlers.NewAuthHandler(repo)

	// Register & Login (public)
	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)

	// Protected route
	authRouter.GET("/profile", middlewares.VerifyToken(rdb), middlewares.Access("user", "admin"), handler.Profile)
	authRouter.DELETE("/logout", middlewares.VerifyToken(rdb), handler.Logout)
	authRouter.PATCH("/update-password", middlewares.VerifyToken(rdb), middlewares.Access("user", "admin"), handler.UpdatePassword)
}
