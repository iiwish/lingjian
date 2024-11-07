package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// TaskService 任务服务
type TaskService struct{}

// CreateScheduledTask 创建定时任务
func (s *TaskService) CreateScheduledTask(appID uint, name, typ, cron string, content map[string]any, timeout, retryTimes int) error {
	// 检查应用是否存在
	var appCount int
	err := model.DB.Get(&appCount, "SELECT COUNT(*) FROM apps WHERE id = ?", appID)
	if err != nil {
		return err
	}
	if appCount == 0 {
		return errors.New("应用不存在")
	}

	// 检查任务名称是否已存在
	var taskCount int
	err = model.DB.Get(&taskCount, "SELECT COUNT(*) FROM scheduled_tasks WHERE app_id = ? AND name = ?", appID, name)
	if err != nil {
		return err
	}
	if taskCount > 0 {
		return errors.New("任务名称已存在")
	}

	// 序列化任务内容
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 创建任务
	_, err = model.DB.Exec(`
		INSERT INTO scheduled_tasks (
			app_id, name, type, cron, content, timeout, retry_times, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, appID, name, typ, cron, string(contentJSON), timeout, retryTimes, 1, time.Now(), time.Now())

	return err
}

// UpdateScheduledTask 更新定时任务
func (s *TaskService) UpdateScheduledTask(taskID uint, name, cron string, content map[string]any, timeout, retryTimes int) error {
	// 检查任务是否存在
	var taskCount int
	err := model.DB.Get(&taskCount, "SELECT COUNT(*) FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		return err
	}
	if taskCount == 0 {
		return errors.New("任务不存在")
	}

	// 序列化任务内容
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 更新任务
	_, err = model.DB.Exec(`
		UPDATE scheduled_tasks
		SET name = ?, cron = ?, content = ?, timeout = ?, retry_times = ?, updated_at = ?
		WHERE id = ?
	`, name, cron, string(contentJSON), timeout, retryTimes, time.Now(), taskID)

	return err
}

// ToggleTaskStatus 切换任务状态
func (s *TaskService) ToggleTaskStatus(taskID uint, status int) error {
	// 检查任务是否存在
	var taskCount int
	err := model.DB.Get(&taskCount, "SELECT COUNT(*) FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		return err
	}
	if taskCount == 0 {
		return errors.New("任务不存在")
	}

	// 更新状态
	_, err = model.DB.Exec(`
		UPDATE scheduled_tasks
		SET status = ?, updated_at = ?
		WHERE id = ?
	`, status, time.Now(), taskID)

	return err
}

// GetTaskLogs 获取任务日志
func (s *TaskService) GetTaskLogs(taskID uint, limit, offset int) ([]map[string]interface{}, error) {
	var logs []map[string]interface{}
	query := `
		SELECT *
		FROM task_logs
		WHERE task_id = ?
		ORDER BY start_time DESC
		LIMIT ? OFFSET ?
	`
	err := model.DB.Select(&logs, query, taskID, limit, offset)
	return logs, err
}

// ExecuteTask 执行任务
func (s *TaskService) ExecuteTask(taskID uint) error {
	// 检查任务是否存在且启用
	var task struct {
		Type    string
		Content string
		Status  int
	}
	err := model.DB.Get(&task, "SELECT type, content, status FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		return errors.New("任务不存在")
	}
	if task.Status != 1 {
		return errors.New("任务未启用")
	}

	// 记录执行开始时间
	startTime := time.Now()

	// TODO: 实现任务执行逻辑
	// 这里需要根据任务类型（SQL/HTTP）实现具体的执行逻辑

	// 记录执行结果
	_, err = model.DB.Exec(`
		INSERT INTO task_logs (task_id, status, result, error, start_time, end_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`, taskID, 1, "执行成功", "", startTime, time.Now())

	return err
}

// CreateElementTrigger 创建元素触发器
func (s *TaskService) CreateElementTrigger(appID uint, elementType string, elementID uint, triggerPoint, typ string, content map[string]any) error {
	// 检查应用是否存在
	var appCount int
	err := model.DB.Get(&appCount, "SELECT COUNT(*) FROM apps WHERE id = ?", appID)
	if err != nil {
		return err
	}
	if appCount == 0 {
		return errors.New("应用不存在")
	}

	// 序列化触发器内容
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 创建触发器
	_, err = model.DB.Exec(`
		INSERT INTO element_triggers (
			app_id, element_type, element_id, trigger_point, type, content, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, appID, elementType, elementID, triggerPoint, typ, string(contentJSON), 1, time.Now(), time.Now())

	return err
}
