package handlers

func (h *BaseHandler)PostMonitor(c *gin.Context) {
	m := models.Monitor()
	if err := c.ShouldBindJSON(m); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}
	if h.Insert(m) != nil {
		utils.Error(c, 500, "数据库错误")
		return
	}
	utils.Success(c, "OK")
}