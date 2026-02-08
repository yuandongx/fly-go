package handlers

import (
	"fly-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func (h *BaseHandler) GetFundList(c *gin.Context) {
	// get fund list from mongoDB
	h.DefaultGetListQuery(collectionFund, c)
}

func (h *BaseHandler) GetFundDetail(c *gin.Context) {
	var query utils.Query
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	utils.Success(c, nil)
}

func (h *BaseHandler) CreateFund(c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) UpdateFund(c *gin.Context) {
	utils.Success(c, nil)
}

func (h *BaseHandler) DeleteFund(c *gin.Context) {
	utils.Success(c, nil)
}
