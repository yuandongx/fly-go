package handlers

import (
	"context"
	"fly-go/database"
	"fly-go/server/models"
	"fly-go/server/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MonitorResource 监控资源处理器
type MonitorResource struct {
	*BaseHandler
	collection string
}

func NewMonitorResource(db *database.MongoDB) *MonitorResource {
	return &MonitorResource{
		BaseHandler: NewBaseHandler(db),
		collection:  "monitor",
	}
}

func (r *MonitorResource) Name() string {
	return "monitor"
}

func (r *MonitorResource) Routes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "", r.List(r.collection), "获取监控列表"},
		{MethodPost, "", r.CreateMonitor(), "创建监控"},
	}
}

func (r *MonitorResource) CustomRoutes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "/:id", r.GetByID(r.collection), "获取监控详情"},
		{MethodPut, "/:id", r.UpdateMonitor(), "更新监控"},
		{MethodDelete, "/:id", r.Delete(r.collection), "删除监控"},
	}
}

// CreateMonitor 创建监控记录
func (r *MonitorResource) CreateMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.Monitor
		if err := c.ShouldBindJSON(&m); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		if _, err := r.Mongo.InsertOne(context.Background(), r.collection, m); err != nil {
			utils.Error(c, 500, "创建失败")
			return
		}
		utils.Success(c, "创建成功")
	}
}

// UpdateMonitor 更新监控记录
func (r *MonitorResource) UpdateMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var m models.Monitor
		if err := c.ShouldBindJSON(&m); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.BadRequest(c, "无效的ID")
			return
		}
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": m}
		if _, err := r.Mongo.UpdateOne(context.Background(), r.collection, filter, update); err != nil {
			utils.Error(c, 500, "更新失败")
			return
		}
		utils.Success(c, "更新成功")
	}
}
