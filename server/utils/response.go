// Package utils provides the utility functions for the application.
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

func ErrorWithData(c *gin.Context, statusCode int, message string, data interface{}) {

	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// BadRequest 返回400错误响应，包装了Error函数
// c: Gin上下文对象
// message: 返回给客户端的错误信息
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 返回401未授权错误响应，使用指定的错误消息
// c: Gin上下文对象
// message: 要返回给客户端的错误消息
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 返回403 Forbidden错误响应，使用指定的错误信息
// c: Gin上下文对象
// message: 返回给客户端的错误信息
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound 返回404状态码的HTTP错误响应，并附带指定的错误信息
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// ServerServerError 返回500状态码和指定的错误信息
// c: Gin上下文对象
// message: 要返回的错误信息
func ServerServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
