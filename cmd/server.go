// Package main provides the entry point for the application.
package main

// main function
import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fly-go/internal/config"
	"fly-go/internal/database"
	"fly-go/internal/routes"
	log "fly-go/logger"

	"github.com/gin-gonic/gin"
)

func Server() {
	logger := log.DefaultLogger()

	logger.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config")
	}

	mongoDB, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", log.Zap("Error", err.Error()))
	}
	defer mongoDB.Close()

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	routes.SetupRoutes(r, mongoDB, logger)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Server starting", log.Zap("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", log.Zap("Error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", log.Zap("Error", err.Error()))
	}

	logger.Info("Server exited")
}
