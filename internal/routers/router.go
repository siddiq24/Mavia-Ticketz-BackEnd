package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	docs "github.com/siddiq24/Tickitz-DB/docs"
	"github.com/siddiq24/Tickitz-DB/internal/middlewares"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swag/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	Ping_Router(r, db)
	InitAuthRouter(r, db)
	InitMovieRouter(r, db)
	InitScheduleRouter(r, db)
	InitSeatRouter(r, db)
	InitOrderRouter(r, db)
	InitProfileRouter(r, db)

	//NOT FOUND
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"status": "404 NOT FOUND"})
	})

	return r
}
