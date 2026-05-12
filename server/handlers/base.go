package handlers

import (
	"context"
	"fly-go/database"
	"fly-go/server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

// BaseHandler 通用基础处理器
type BaseHandler struct {
	Mongo *database.MongoDB
}

func NewBaseHandler(mongoDB *database.MongoDB) *BaseHandler {
	return &BaseHandler{Mongo: mongoDB}
}

// ==================== 通用 CRUD 操作 ====================

// List 查询列表
func (h *BaseHandler) List(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := utils.BaseQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		total, items := h.Mongo.Find(c.Request.Context(), collection, query)
		utils.Success(c, gin.H{"total": total, "items": items})
	}
}

// GetByID 根据ID查询单条
func (h *BaseHandler) GetByID(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.BadRequest(c, "ID不能为空")
			return
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.BadRequest(c, "无效的ID格式")
			return
		}
		result := h.Mongo.FindOne(c.Request.Context(), collection, bson.M{"_id": objectID})
		var doc bson.M
		if err := result.Decode(&doc); err != nil {
			if err == mongo.ErrNoDocuments {
				utils.NotFound(c, "记录不存在")
			} else {
				utils.Error(c, 500, "查询失败")
			}
			return
		}
		utils.Success(c, doc)
	}
}

// Create 创建记录
func (h *BaseHandler) Create(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var obj bson.M
		if err := c.ShouldBindJSON(&obj); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		if _, err := h.Mongo.InsertOne(context.Background(), collection, obj); err != nil {
			utils.Error(c, 500, "创建失败")
			return
		}
		utils.Success(c, "创建成功")
	}
}

// Update 更新记录
func (h *BaseHandler) Update(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.BadRequest(c, "ID不能为空")
			return
		}
		var obj bson.M
		if err := c.ShouldBindJSON(&obj); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.BadRequest(c, "无效的ID格式")
			return
		}
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": obj}
		_, err = h.Mongo.UpdateOne(context.Background(), collection, filter, update)
		if err != nil {
			utils.Error(c, 500, "更新失败")
			return
		}
		utils.Success(c, "更新成功")
	}
}

// Delete 删除记录
func (h *BaseHandler) Delete(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.BadRequest(c, "ID不能为空")
			return
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.BadRequest(c, "无效的ID格式")
			return
		}
		_, err = h.Mongo.DeleteOne(c.Request.Context(), collection, bson.M{"_id": objectID})
		if err != nil {
			utils.Error(c, 500, "删除失败")
			return
		}
		utils.Success(c, "删除成功")
	}
}

// ==================== 健康检查 ====================

func HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}

// ==================== Resource 接口 (用于自动注册) ====================

// RouteConfig 路由配置
type RouteConfig struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Description string
}

// Resource 资源接口 - 实现此接口即可自动注册路由
type Resource interface {
	// Name 返回资源名称（用于路由路径）
	Name() string
	// Routes 返回该资源的所有路由配置
	Routes() []RouteConfig
	// CustomRoutes 返回自定义路由（用于特殊业务逻辑）
	CustomRoutes() []RouteConfig
}

// RegisterResources 注册所有实现了Resource接口的资源
func RegisterResources(r *gin.RouterGroup, resources []Resource, mongoDB *database.MongoDB, logger *zap.Logger) {
	for _, resource := range resources {
		name := resource.Name()
		prefix := strings.ToLower(name)

		// 注册标准 CRUD 路由
		routes := resource.Routes()
		for _, route := range routes {
			fullPath := "/" + prefix + route.Path
			r.Handle(route.Method, fullPath, route.Handler)
			logger.Info("Route registered",
				zap.String("method", route.Method),
				zap.String("path", fullPath),
				zap.String("description", route.Description),
			)
		}

		// 注册自定义路由
		customRoutes := resource.CustomRoutes()
		for _, route := range customRoutes {
			fullPath := "/" + prefix + route.Path
			r.Handle(route.Method, fullPath, route.Handler)
			logger.Info("Custom route registered",
				zap.String("method", route.Method),
				zap.String("path", fullPath),
			)
		}
	}
}
