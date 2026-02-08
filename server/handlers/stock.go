package handlers

import (
	"fly-go/server/utils"

	"github.com/gin-gonic/gin"
)

func (h *BaseHandler) GetStockList(c *gin.Context) {
	h.DefaultGetListQuery(h.collection, c)
}

func (h *BaseHandler) GetStockDetail(c *gin.Context) {
	query := utils.Query{}
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	utils.Success(c, nil)
}

func (h *BaseHandler) CreateStock(c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) UpdateStock(c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) DeleteStock(c *gin.Context) {
	utils.Success(c, nil)
}
