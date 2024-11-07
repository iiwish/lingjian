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
	AppID      uint           `json:"app_id" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Type       string         `json:"type" binding:"required"`
	Cron       string         `json:"cron" binding:"required"`
	Content    map[string]any `json:"content" binding:"required"`
	Timeout    int            `json:"timeout"`
	RetryTimes int            `json:"retry_times"`
}

// @Summary      创建定时任务
// @Description  创建新的定时任务
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        request body CreateScheduledTaskRequest true "创建定时任务请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/scheduled [post]
func CreateScheduledTask(c *gin.Context) {
	var req CreateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	taskService := &service.TaskService{}
	if err := taskService.CreateScheduledTask(
		req.AppID,
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

// @Summary      更新定时任务
// @Description  更新已存在的定时任务
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id path int true "任务ID"
// @Param        request body UpdateScheduledTaskRequest true "更新定时任务请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/scheduled/{id} [put]
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

// @Summary      切换任务状态
// @Description  启用或禁用定时任务
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id path int true "任务ID"
// @Param        status query int true "状态码"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/scheduled/{id}/toggle [post]
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

// @Summary      获取任务日志
// @Description  获取定时任务的执行日志
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id path int true "任务ID"
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/scheduled/{id}/logs [get]
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

// @Summary      执行任务
// @Description  手动执行定时任务
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id path int true "任务ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/scheduled/{id}/execute [post]
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
	AppID        uint           `json:"app_id" binding:"required"`
	ElementType  string         `json:"element_type" binding:"required"`
	ElementID    uint           `json:"element_id" binding:"required"`
	TriggerPoint string         `json:"trigger_point" binding:"required"`
	Type         string         `json:"type" binding:"required"`
	Content      map[string]any `json:"content" binding:"required"`
}

// @Summary      创建元素触发器
// @Description  创建新的元素触发器
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        request body CreateElementTriggerRequest true "创建元素触发器请求参数"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /tasks/triggers [post]
func CreateElementTrigger(c *gin.Context) {
	var req CreateElementTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "无效的请求参数")
		return
	}

	taskService := &service.TaskService{}
	if err := taskService.CreateElementTrigger(
		req.AppID,
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
