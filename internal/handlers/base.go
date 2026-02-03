package handlers

import "fly-go/internal/database"

type BaseHandler struct {
	mongoDB *database.MongoDB
}

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BaseQuery struct {
	Page    int    `json:"page" form:"page" validate:"required,min=1" default:"1"`
	Size    int    `json:"size" form:"size" validate:"required,min=1,max=100" default:"10"`
	Search  string `json:"search" form:"search" validate:"omitempty"`
	OrderBy string `json:"order_by" form:"order_by" validate:"omitempty"`
	Order   string `json:"order" form:"order" validate:"omitempty"`
}

type Query struct {
	ID string `json:"id" form:"id" validate:"required"`
}

func NewBaseHandler(mongoDB *database.MongoDB) *BaseHandler {
	return &BaseHandler{
		mongoDB: mongoDB,
	}
}

func (h *BaseHandler) GetMongoDB() *database.MongoDB {
	return h.mongoDB
}
