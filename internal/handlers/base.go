package handlers

import (
	"fly-go/internal/database"
	"fly-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	Mongo *database.MongoDB
}

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewBaseHandler(mongoDB *database.MongoDB) *BaseHandler {
	return &BaseHandler{
		Mongo: mongoDB,
	}
}

func (h *BaseHandler) GetMongoDB() *database.MongoDB {
	return h.Mongo
}

func (h *BaseHandler) DefaultGetListQuery(collection string, c *gin.Context) {
	query := &utils.BaseQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		utils.Error(c, 400, "参数错误")
	} else {
		results := h.Mongo.Find(c.Request.Context(), collection, *query)
		utils.Success(c, results)
	}
}

const collectionFund = "fund"
const collectionStock = "stock"
const collectionTask = "task"

// const collectionFollow = "follow"
