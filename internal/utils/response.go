// Package utils provides the utility functions for the application.
package utils

import (
	"net/http"

	"fly-go/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	logger.Error("API Error",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("message", message))

	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

func ErrorWithData(c *gin.Context, statusCode int, message string, data interface{}) {
	logger.Error("API Error",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("message", message))

	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
