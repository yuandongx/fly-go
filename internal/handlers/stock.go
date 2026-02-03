package handlers

import (
	"context"
	"fly-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func (h *BaseHandler) GetStockList(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) GetStockDetail(ctx context.Context, c *gin.Context) {
	var query Query
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	utils.Success(c, nil)
}

func (h *BaseHandler) CreateStock(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) UpdateStock(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) DeleteStock(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}
