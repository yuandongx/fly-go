package handlers

import (
	"context"
	"fly-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func (h *BaseHandler) GetFundList(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) GetFundDetail(ctx context.Context, c *gin.Context) {
	var query Query
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	utils.Success(c, nil)
}

func (h *BaseHandler) CreateFund(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) UpdateFund(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) DeleteFund(ctx context.Context, c *gin.Context) {
	utils.Success(c, nil)
}
