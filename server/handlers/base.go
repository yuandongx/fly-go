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

// DefaultGetListQuery 处理默认的列表查询请求，从指定集合中查询数据并返回分页结果
func (h *BaseHandler) DefaultGetListQuery(collection string, c *gin.Context) {
	query := &utils.BaseQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		utils.Error(c, 400, "参数错误")
	} else {
		total, results := h.Mongo.Find(c.Request.Context(), collection, *query)
		utils.Success(c, gin.H{"total": total, "items": results})
	}
}

// Insert 向数据库集合中插入一个文档
func (h *BaseHandler) Insert(obj interface{}) error {
	_, err := h.Mongo.InsertOne(context.Background(), h.collection, obj)
	return err
}

// UpdateByID 根据ID更新文档
// 参数:
//   - id: 文档ID
//   - obj: 要更新的对象
//
// 返回:
//   - error: 更新失败时返回错误
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
