// Package main provides the entry point for the application.
package main

// main function
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fly-go/internal/config"
	"fly-go/internal/database"
	"fly-go/internal/routes"
	"fly-go/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	if err := database.Connect(cfg.Database); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Disconnect()

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	routes.SetupRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
