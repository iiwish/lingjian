package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/utils"
)

// RegisterTaskRoutes 注册任务相关路由
func RegisterTaskRoutes(r *gin.RouterGroup) {
	task := r.Group("/tasks")
	{
		// 定时任务
		task.POST("/scheduled", CreateScheduledTask)
		task.PUT("/scheduled/:id", UpdateScheduledTask)
		task.POST("/scheduled/:id/toggle", ToggleTaskStatus)
		task.GET("/scheduled/:id/logs", GetTaskLogs)
		task.POST("/scheduled/:id/execute", ExecuteTask)

		// 元素触发器
		task.POST("/triggers", CreateElementTrigger)
	}
}

type CreateScheduledTaskRequest struct {
	ApplicationID uint           `json:"application_id" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	Type          string         `json:"type" binding:"required"`
	Cron          string         `json:"cron" binding:"required"`
	Content       map[string]any `json:"content" binding:"required"`
	Timeout       int            `json:"timeout"`
	RetryTimes    int            `json:"retry_times"`
}

// CreateScheduledTask 创建定时任务
func CreateScheduledTask(c *gin.Context) {
	var req CreateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	taskService := &service.TaskService{}
	if err := taskService.CreateScheduledTask(
		req.ApplicationID,
		req.Name,
		req.Type,
		req.Cron,
		req.Content,
		req.Timeout,
		req.RetryTimes,
	); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

type UpdateScheduledTaskRequest struct {
	Name       string         `json:"name" binding:"required"`
	Cron       string         `json:"cron" binding:"required"`
	Content    map[string]any `json:"content" binding:"required"`
	Timeout    int            `json:"timeout"`
	RetryTimes int            `json:"retry_times"`
}

// UpdateScheduledTask 更新定时任务
func UpdateScheduledTask(c *gin.Context) {
	var req UpdateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	taskID := utils.ParseUint(c.Param("id"))
	taskService := &service.TaskService{}
	if err := taskService.UpdateScheduledTask(
		taskID,
		req.Name,
		req.Cron,
		req.Content,
		req.Timeout,
		req.RetryTimes,
	); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// ToggleTaskStatus 切换任务状态
func ToggleTaskStatus(c *gin.Context) {
	taskID := utils.ParseUint(c.Param("id"))
	status := utils.ParseInt(c.Query("status"))

	taskService := &service.TaskService{}
	if err := taskService.ToggleTaskStatus(taskID, status); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetTaskLogs 获取任务执行日志
func GetTaskLogs(c *gin.Context) {
	taskID := utils.ParseUint(c.Param("id"))
	limit := utils.ParseInt(c.DefaultQuery("limit", "10"))
	offset := utils.ParseInt(c.DefaultQuery("offset", "0"))

	taskService := &service.TaskService{}
	logs, err := taskService.GetTaskLogs(taskID, limit, offset)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, logs)
}

// ExecuteTask 手动执行任务
func ExecuteTask(c *gin.Context) {
	taskID := utils.ParseUint(c.Param("id"))

	taskService := &service.TaskService{}
	if err := taskService.ExecuteTask(taskID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateElementTriggerRequest struct {
	ApplicationID uint           `json:"application_id" binding:"required"`
	ElementType   string         `json:"element_type" binding:"required"`
	ElementID     uint           `json:"element_id" binding:"required"`
	TriggerPoint  string         `json:"trigger_point" binding:"required"`
	Type          string         `json:"type" binding:"required"`
	Content       map[string]any `json:"content" binding:"required"`
}

// CreateElementTrigger 创建元素触发器
func CreateElementTrigger(c *gin.Context) {
	var req CreateElementTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	taskService := &service.TaskService{}
	if err := taskService.CreateElementTrigger(
		req.ApplicationID,
		req.ElementType,
		req.ElementID,
		req.TriggerPoint,
		req.Type,
		req.Content,
	); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}
