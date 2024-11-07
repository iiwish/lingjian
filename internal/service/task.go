package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// 任务类型常量
const (
	TaskTypeSQL  = "sql"
	TaskTypeHTTP = "http"
)

// 任务状态常量
const (
	TaskStatusDisabled = 0
	TaskStatusEnabled  = 1
	TaskStatusRunning  = 2
)

// 触发点常量
const (
	TriggerPointBefore = "before"
	TriggerPointAfter  = "after"
)

// TaskService 任务服务
type TaskService struct {
	runningTasks sync.Map
}

// CreateScheduledTask 创建定时任务
func (s *TaskService) CreateScheduledTask(appID uint, name, typ, cron string, content map[string]interface{}, timeout, retryTimes int) error {
	// 检查任务类型
	if typ != TaskTypeSQL && typ != TaskTypeHTTP {
		return fmt.Errorf("不支持的任务类型: %s", typ)
	}

	// 检查应用是否存在
	var appCount int
	err := model.DB.Get(&appCount, "SELECT COUNT(*) FROM apps WHERE id = ?", appID)
	if err != nil {
		return fmt.Errorf("检查应用失败: %v", err)
	}
	if appCount == 0 {
		return errors.New("应用不存在")
	}

	// 检查任务名称是否已存在
	var taskCount int
	err = model.DB.Get(&taskCount, "SELECT COUNT(*) FROM scheduled_tasks WHERE app_id = ? AND name = ?", appID, name)
	if err != nil {
		return fmt.Errorf("检查任务名称失败: %v", err)
	}
	if taskCount > 0 {
		return errors.New("任务名称已存在")
	}

	// 验证任务内容
	if err := s.validateTaskContent(typ, content); err != nil {
		return fmt.Errorf("验证任务内容失败: %v", err)
	}

	// 序列化任务内容
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("序列化任务内容失败: %v", err)
	}

	// 设置默认值
	if timeout <= 0 {
		timeout = 60 // 默认60秒超时
	}
	if retryTimes < 0 {
		retryTimes = 0
	}

	// 创建任务
	_, err = model.DB.Exec(`
		INSERT INTO scheduled_tasks (
			app_id, name, type, cron, content, timeout, retry_times, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, appID, name, typ, cron, string(contentJSON), timeout, retryTimes, TaskStatusEnabled, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("创建任务失败: %v", err)
	}

	return nil
}

// UpdateScheduledTask 更新定时任务
func (s *TaskService) UpdateScheduledTask(taskID uint, name, cron string, content map[string]interface{}, timeout, retryTimes int) error {
	// 检查任务是否存在
	var task struct {
		Type   string
		Status int
	}
	err := model.DB.Get(&task, "SELECT type, status FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("任务不存在")
		}
		return err
	}

	// 检查任务是否正在运行
	if task.Status == TaskStatusRunning {
		return errors.New("任务正在运行中，无法更新")
	}

	// 验证任务内容
	if err := s.validateTaskContent(task.Type, content); err != nil {
		return err
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
	// 检查状态值是否有效
	if status != TaskStatusEnabled && status != TaskStatusDisabled {
		return errors.New("无效的状态值")
	}

	// 检查任务是否存在
	var currentStatus int
	err := model.DB.Get(&currentStatus, "SELECT status FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("任务不存在")
		}
		return err
	}

	// 检查任务是否正在运行
	if currentStatus == TaskStatusRunning {
		return errors.New("任务正在运行中，无法更改状态")
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
		ID         uint
		Type       string
		Content    string
		Status     int
		Timeout    int
		RetryTimes int
	}
	err := model.DB.Get(&task, "SELECT id, type, content, status, timeout, retry_times FROM scheduled_tasks WHERE id = ?", taskID)
	if err != nil {
		return errors.New("任务不存在")
	}
	if task.Status != TaskStatusEnabled {
		return errors.New("任务未启用")
	}

	// 检查任务是否已在运行
	if _, running := s.runningTasks.LoadOrStore(taskID, true); running {
		return errors.New("任务正在运行中")
	}
	defer s.runningTasks.Delete(taskID)

	// 更新任务状态为运行中
	_, err = model.DB.Exec("UPDATE scheduled_tasks SET status = ? WHERE id = ?", TaskStatusRunning, taskID)
	if err != nil {
		return err
	}
	defer model.DB.Exec("UPDATE scheduled_tasks SET status = ? WHERE id = ?", TaskStatusEnabled, taskID)

	// 记录执行开始时间
	startTime := time.Now()
	var result string
	var execErr error

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.Timeout)*time.Second)
	defer cancel()

	// 解析任务内容
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(task.Content), &content); err != nil {
		return err
	}

	// 执行任务（带重试）
	for i := 0; i <= task.RetryTimes; i++ {
		select {
		case <-ctx.Done():
			execErr = ctx.Err()
			break
		default:
			// 根据任务类型执行
			switch task.Type {
			case TaskTypeSQL:
				result, execErr = s.executeSQL(content)
			case TaskTypeHTTP:
				result, execErr = s.executeHTTP(content)
			default:
				execErr = errors.New("不支持的任务类型")
			}

			if execErr == nil {
				break
			}

			// 最后一次重试失败
			if i == task.RetryTimes {
				break
			}

			// 等待一段时间后重试
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	// 记录执行结果
	status := 1
	if execErr != nil {
		status = 0
		result = execErr.Error()
	}

	_, err = model.DB.Exec(`
		INSERT INTO task_logs (task_id, status, result, error, start_time, end_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`, taskID, status, result, execErr, startTime, time.Now())

	return execErr
}

// executeSQL 执行SQL任务
func (s *TaskService) executeSQL(content map[string]interface{}) (string, error) {
	sql, ok := content["sql"].(string)
	if !ok {
		return "", errors.New("无效的SQL语句")
	}

	// SQL安全检查
	if err := s.validateSQL(sql); err != nil {
		return "", err
	}

	// 执行SQL
	result, err := model.DB.Exec(sql)
	if err != nil {
		return "", err
	}

	// 获取影响行数
	affected, _ := result.RowsAffected()
	return fmt.Sprintf("执行成功，影响 %d 行", affected), nil
}

// executeHTTP 执行HTTP任务
func (s *TaskService) executeHTTP(content map[string]interface{}) (string, error) {
	url, ok := content["url"].(string)
	if !ok {
		return "", errors.New("无效的URL")
	}

	method, _ := content["method"].(string)
	if method == "" {
		method = "GET"
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	// 添加请求头
	if headers, ok := content["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Add(key, strValue)
			}
		}
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP请求失败: %d %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// validateTaskContent 验证任务内容
func (s *TaskService) validateTaskContent(typ string, content map[string]interface{}) error {
	switch typ {
	case TaskTypeSQL:
		sql, ok := content["sql"].(string)
		if !ok {
			return errors.New("SQL任务必须包含sql字段")
		}
		return s.validateSQL(sql)

	case TaskTypeHTTP:
		url, ok := content["url"].(string)
		if !ok {
			return errors.New("HTTP任务必须包含url字段")
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return errors.New("无效的URL格式")
		}
		return nil

	default:
		return errors.New("不支持的任务类型")
	}
}

// validateSQL SQL安全检查
func (s *TaskService) validateSQL(sql string) error {
	sqlUpper := strings.ToUpper(sql)

	// 禁止危险操作
	dangerousKeywords := []string{
		"DROP", "TRUNCATE", "ALTER", "CREATE",
		"GRANT", "REVOKE", "RENAME",
	}

	// 检查每个关键字是否作为独立的单词出现
	for _, keyword := range dangerousKeywords {
		// 在关键字前后添加空格，以确保匹配完整的单词
		pattern := " " + keyword + " "
		if strings.Contains(" "+sqlUpper+" ", pattern) {
			return fmt.Errorf("SQL语句包含危险关键字: %s", keyword)
		}
	}

	return nil
}

// CreateElementTrigger 创建元素触发器
func (s *TaskService) CreateElementTrigger(appID uint, elementType string, elementID uint, triggerPoint, typ string, content map[string]interface{}) error {
	// 检查触发点是否有效
	if triggerPoint != TriggerPointBefore && triggerPoint != TriggerPointAfter {
		return errors.New("无效的触发点")
	}

	// 检查应用是否存在
	var appCount int
	err := model.DB.Get(&appCount, "SELECT COUNT(*) FROM apps WHERE id = ?", appID)
	if err != nil {
		return err
	}
	if appCount == 0 {
		return errors.New("应用不存在")
	}

	// 验证触发器内容
	if err := s.validateTaskContent(typ, content); err != nil {
		return err
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
	`, appID, elementType, elementID, triggerPoint, typ, string(contentJSON), TaskStatusEnabled, time.Now(), time.Now())

	return err
}
