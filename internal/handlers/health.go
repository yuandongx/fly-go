// Package handlers provides the HTTP request handlers for the application.
package handlers

import (
	"fly-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}
