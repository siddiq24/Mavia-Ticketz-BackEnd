package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
)

func InitDashboardRouter(router *gin.Engine, db *pgxpool.Pool) {
	repo := repositories.NewDashboardRepository(db)
	handler := handlers.NewDashboardHandler(repo)

	admin := router.Group("/admin")
	{
		admin.GET("/dashboard", handler.GetDashboard)
	}
}
