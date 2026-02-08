// Package server provides the server implementation for the application.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fly-go/config"
	"fly-go/database"
	log "fly-go/logger"
	"fly-go/server/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Start 启动应用程序服务器并监听指定端口
//
// 参数:
//
//	port: 服务器监听的端口号
//
// 功能:
//   - 加载应用配置
//   - 初始化数据库连接
//   - 设置Gin路由
//   - 启动HTTP服务器
//   - 处理优雅关机信号
//
// 注意:
//   - 使用goroutine异步启动服务器
//   - 捕获SIGINT和SIGTERM信号进行优雅关机
//   - 数据库连接会在函数结束时自动关闭
func Start(port int) {
	logger := log.DefaultLogger()

	logger.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config")
	}

	mongoDB, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		logger.Info("Mongo connect: ", zap.String("Host", cfg.Database.MongoHost), zap.Int("Port", cfg.Database.MongoPort))
		logger.Error("Failed to connect to database", zap.String("Error", err.Error()))
	}

	defer mongoDB.Close()

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	routes.SetupRoutes(r, mongoDB, logger)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		logger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", zap.String("Error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.String("Error", err.Error()))
	}

	logger.Info("Server exited")
}
