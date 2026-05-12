package routes

import (
	"fly-go/database"
	log "fly-go/logger"
	"fly-go/server/handlers"
	"fly-go/server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, mongoDB *database.MongoDB, logger *log.ILogger) {
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", handlers.HealthCheck)

		// 注册所有资源路由
		resources := []handlers.Resource{
			handlers.NewStockResource(mongoDB),
			handlers.NewFundResource(mongoDB),
			handlers.NewTaskResource(mongoDB),
			handlers.NewMonitorResource(mongoDB),
		}
		handlers.RegisterResources(api, resources, mongoDB, logger.Logger())
	}
}
