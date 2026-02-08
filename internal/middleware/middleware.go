// Package middleware provides the middleware functions for the application.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	log "fly-go/logger"
)

func Logger(logger *log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request",
			log.Zap("method", method),
			log.Zap("path", path),
			log.Zap("status", statusCode),
			log.Zap("latency", latency),
			log.Zap("client_ip", c.ClientIP()),
		)
	}
}

func Recovery(logger *log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					log.Zap("path", c.Request.URL.Path),
					log.Zap("method", c.Request.Method),
					log.Zap("error", err),
				)
				c.JSON(500, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
