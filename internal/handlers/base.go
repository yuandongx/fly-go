package handlers

import "fly-go/internal/database"

type BaseHandler struct {
	mongoDB *database.MongoDB
}

func NewBaseHandler(mongoDB *database.MongoDB) *BaseHandler {
	return &BaseHandler{
		mongoDB: mongoDB,
	}
}

func (h *BaseHandler) GetMongoDB() *database.MongoDB {
	return h.mongoDB
}
