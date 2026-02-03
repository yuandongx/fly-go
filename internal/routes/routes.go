// Package routes provides the HTTP routes configuration for the application.
package routes

import (
	"fly-go/internal/database"
	"fly-go/internal/handlers"
	"fly-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, mongoDB *database.MongoDB) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	baseHandler := handlers.NewBaseHandler(mongoDB)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/health", baseHandler.Check)
			v1.GET("/stock/list", baseHandler.GetStockList)
			v1.GET("/fund/list", baseHandler.GetFundList)
			v1.GET("/task/list", baseHandler.GetTaskList)
		}
	}
}
