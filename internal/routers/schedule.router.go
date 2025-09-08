package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool) {
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)

	scheduleRouter := router.Group("/schedules")
	{
		scheduleRouter.GET("", scheduleHandler.GetSchedules)
	}
}
