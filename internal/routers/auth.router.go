package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	repo := repositories.NewUserRepository(db)
	handler := handlers.NewAuthHandler(repo)

	// Register & Login (public)
	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)

	// Protected route
	authRouter.GET("/profile", middlewares.VerifyToken, middlewares.Access("user", "admin"), handler.Profile)
}
