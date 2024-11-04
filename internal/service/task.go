package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/queue"
)

type TaskService struct{}

// TaskMessage 任务消息结构
type TaskMessage struct {
	TaskID  uint           `json:"task_id"`
	Type    string         `json:"type"`
	Content map[string]any `json:"content"`
	Timeout int            `json:"timeout"`
}

// CreateScheduledTask 创建定时任务
func (s *TaskService) CreateScheduledTask(appID uint, name, taskType, cron string, content map[string]any, timeout, retryTimes int) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = model.DB.Exec(`
		INSERT INTO scheduled_tasks (
			app_id, name, type, cron,
			content, timeout, retry_times, status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, appID, name, taskType, cron,
		contentJSON, timeout, retryTimes,
		time.Now(), time.Now())

	return err
}

// UpdateScheduledTask 更新定时任务
func (s *TaskService) UpdateScheduledTask(taskID uint, name, cron string, content map[string]any, timeout, retryTimes int) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	result, err := model.DB.Exec(`
		UPDATE scheduled_tasks SET
		name = ?, cron = ?, content = ?,
		timeout = ?, retry_times = ?,
		updated_at = ?
		WHERE id = ?
	`, name, cron, contentJSON,
		timeout, retryTimes,
		time.Now(), taskID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("任务不存在")
	}

	return nil
}

// ToggleTaskStatus 切换任务状态
func (s *TaskService) ToggleTaskStatus(taskID uint, status int) error {
	result, err := model.DB.Exec(`
		UPDATE scheduled_tasks SET
		status = ?, updated_at = ?
		WHERE id = ?
	`, status, time.Now(), taskID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("任务不存在")
	}

	return nil
}

// GetTaskLogs 获取任务执行日志
func (s *TaskService) GetTaskLogs(taskID uint, limit, offset int) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	err := model.DB.Select(&logs, `
		SELECT * FROM task_logs
		WHERE task_id = ?
		ORDER BY start_time DESC
		LIMIT ? OFFSET ?
	`, taskID, limit, offset)
	return logs, err
}

// CreateElementTrigger 创建元素触发器
func (s *TaskService) CreateElementTrigger(
	appID uint,
	elementType string,
	elementID uint,
	triggerPoint string,
	taskType string,
	content map[string]any,
) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = model.DB.Exec(`
		INSERT INTO element_triggers (
			app_id, element_type, element_id,
			trigger_point, type, content, status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, appID, elementType, elementID,
		triggerPoint, taskType, contentJSON,
		time.Now(), time.Now())

	return err
}

// ExecuteTask 执行任务
func (s *TaskService) ExecuteTask(taskID uint) error {
	var task model.ScheduledTask
	err := model.DB.Get(&task, "SELECT * FROM scheduled_tasks WHERE id = ?", taskID)
	if err == sql.ErrNoRows {
		return errors.New("任务不存在")
	}
	if err != nil {
		return err
	}

	if task.Status != 1 {
		return errors.New("任务已禁用")
	}

	var content map[string]any
	err = json.Unmarshal([]byte(task.Content), &content)
	if err != nil {
		return err
	}

	// 创建任务消息
	message := TaskMessage{
		TaskID:  task.ID,
		Type:    task.Type,
		Content: content,
		Timeout: task.Timeout,
	}

	// 发布到消息队列
	return queue.PublishTask(message)
}

// LogTaskExecution 记录任务执行日志
func (s *TaskService) LogTaskExecution(
	taskID uint,
	status int,
	result string,
	errorMsg string,
	startTime, endTime time.Time,
) error {
	_, err := model.DB.Exec(`
		INSERT INTO task_logs (
			task_id, status, result, error,
			start_time, end_time
		) VALUES (?, ?, ?, ?, ?, ?)
	`, taskID, status, result, errorMsg,
		startTime, endTime)

	return err
}
