package handlers

import (
	"fly-go/database"
	"fly-go/fly"
	"fly-go/server/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskResource 任务资源处理器
type TaskResource struct {
	*BaseHandler
	collection string
}

func NewTaskResource(db *database.MongoDB) *TaskResource {
	return &TaskResource{
		BaseHandler: NewBaseHandler(db),
		collection:  "tasks",
	}
}

func (r *TaskResource) Name() string {
	return "task"
}

func (r *TaskResource) Routes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "", r.List(r.collection), "获取任务列表"},
		{MethodPost, "", r.CreateTask(), "创建任务"},
	}
}

func (r *TaskResource) CustomRoutes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "/:id", r.GetByID(r.collection), "获取任务详情"},
		{MethodPut, "/:id", r.UpdateTask(), "更新任务"},
		{MethodDelete, "/:id", r.Delete(r.collection), "删除任务"},
	}
}

func (r *TaskResource) CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task fly.Runner
		if err := c.ShouldBindJSON(&task); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		if _, err := r.Mongo.InsertOne(c.Request.Context(), r.collection, task); err != nil {
			utils.Error(c, 500, "创建失败")
			return
		}
		utils.Success(c, "创建成功")
	}
}

func (r *TaskResource) UpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task fly.Runner
		if err := c.ShouldBindJSON(&task); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.BadRequest(c, "无效的ID")
			return
		}
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": task}
		if _, err := r.Mongo.UpdateOne(c.Request.Context(), r.collection, filter, update); err != nil {
			utils.Error(c, 500, "更新失败")
			return
		}
		utils.Success(c, "更新成功")
	}
}
