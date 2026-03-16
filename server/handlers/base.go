package handlers

import (
	"context"
	"fly-go/database"
	"fly-go/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	Mongo      *database.MongoDB
	collection string
}

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BM = bson.M

func NewBaseHandler(name string, mongoDB *database.MongoDB) *BaseHandler {
	return &BaseHandler{
		collection: name,
		Mongo:      mongoDB,
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
		total, results := h.Mongo.Find(c.Request.Context(), collection, *query)
		utils.Success(c, gin.H{"total": total, "data": results})
	}
}
func (h *BaseHandler) Insert(obj interface{}) error {
	_, err := h.Mongo.InsertOne(context.Background(), h.collection, obj)
	return err
}

func (h *BaseHandler) UpdateByID(id string, obj interface{}) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_filter := BM{"_id": objectID}
	doc := BM{"$set": obj}
	_, err = h.Mongo.UpdateOne(context.Background(), h.collection, _filter, doc)
	return err
}

func (h *BaseHandler) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	_, err = h.Mongo.DeleteOne(c.Request.Context(), h.collection, BM{"_id": _id})
	if err != nil {
		utils.Error(c, 500, "数据库错误")
		return
	}
	utils.Success(c, "ok")
}
