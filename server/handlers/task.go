package handlers

import (
	"github.com/gin-gonic/gin"

	"fly-go/fly"
	"fly-go/server/utils"
)

// GetTaskList 获取任务列表
func (h *BaseHandler) GetTaskList(c *gin.Context) {
	h.DefaultGetListQuery(h.collection, c)
}

// PostTask 创建新任务
func (h *BaseHandler) PostTask(c *gin.Context) {
	runner := &fly.Runner{}
	if err := c.ShouldBindJSON(runner); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	if h.Insert(runner) != nil {
		utils.Error(c, 500, "数据库错误")
		return
	}
	utils.Success(c, "OK")
}

// UpdateTask 更新任务
func (h *BaseHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	runner := &fly.Runner{}
	if err := c.ShouldBindJSON(runner); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	if h.UpdateByID(id, runner) != nil {
		utils.Error(c, 500, "数据库错误")
		return
	}
	utils.Success(c, "OK")
}

// DeleteTask 删除任务
func (h *BaseHandler) DeleteTask(c *gin.Context) {
	h.DeleteByID(c)
}
