package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/handlers"
)

func Ping_Router(r *gin.Engine, db *pgxpool.Pool) {
	pingR := r.Group("/ping")
	pingH := handlers.NewPingHandler()

	pingR.GET("", pingH.GetPing)
}
