package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRepo := repositories.NewSeatRepository(db)
	seatHandler := handlers.NewSeatHandler(seatRepo)

	seatRouter := router.Group("/seats")
	{
		seatRouter.GET("/available/:scheduleID", seatHandler.GetAvailableSeats)
	}
}
