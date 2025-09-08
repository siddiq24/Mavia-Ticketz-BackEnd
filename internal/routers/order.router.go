package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	orderRouter := router.Group("/orders")
	orderRouter.Use(middlewares.VerifyToken)

	orderRouter.POST("", orderHandler.CreateOrder)
	orderRouter.GET("/history", orderHandler.GetHistory)
}
