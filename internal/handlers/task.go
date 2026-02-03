package handlers

import (
	"github.com/gin-gonic/gin"

	"fly-go/internal/utils"
)

// GetTaskList 获取任务列表
func (h *BaseHandler) GetTaskList(c *gin.Context) {
	h.DefaultGetListQuery(collectionTask, c)
}

// GetTaskDetail 获取任务详情
func (h *BaseHandler) GetTaskDetail(c *gin.Context) {
	query := utils.Query{}
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, 400, "参数错误")
	}
	utils.Success(c, nil)
}

// CreateTask 创建任务
func (h *BaseHandler) CreateTask(c *gin.Context) {
	utils.Success(c, nil)
}

// UpdateTask 更新任务
func (h *BaseHandler) UpdateTask(c *gin.Context) {
	utils.Success(c, nil)
}

// DeleteTask 删除任务
func (h *BaseHandler) DeleteTask(c *gin.Context) {
	utils.Success(c, nil)
}
