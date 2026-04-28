// Package routes provides the HTTP routes configuration for the application.
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

	baseHandler := handlers.NewBaseHandler("default", mongoDB)
	taskHandler := handlers.NewBaseHandler("task", mongoDB)
	stockHandler := handlers.NewBaseHandler("stock", mongoDB)
	fundHandler := handlers.NewBaseHandler("fund", mongoDB)
	montitoHandler := handlers.NewBaseHandler("monitor", mongoDB)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/health", baseHandler.Check)

			v1.GET("/stock", stockHandler.GetStockList)

			v1.GET("/fund", fundHandler.GetFundList)

			// Task routes
			v1.GET("/task", taskHandler.GetTaskList)
			v1.POST("/task", taskHandler.PostTask)
			v1.PUT("/task/:id", taskHandler.UpdateTask)
			v1.DELETE("/task/:id", taskHandler.DeleteTask)

			// Monitor routes
			v1.GET("/monitor", montitoHandler.GetMonitorList)

		}
	}
}
