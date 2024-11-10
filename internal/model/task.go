package model

import "time"

// TaskMessage 任务消息结构
type TaskMessage struct {
	TaskID    uint      `json:"task_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
}

// ScheduledTask 定时任务表结构
type ScheduledTask struct {
	ID         uint      `db:"id" json:"id"`
	AppID      uint      `db:"app_id" json:"app_id"`
	Name       string    `db:"name" json:"name"`
	Type       string    `db:"type" json:"type"`               // sql:SQL任务 http:HTTP任务
	Cron       string    `db:"cron" json:"cron"`               // cron表达式
	Content    string    `db:"content" json:"content"`         // 任务内容（SQL语句或HTTP配置）
	Timeout    int       `db:"timeout" json:"timeout"`         // 超时时间（秒）
	RetryTimes int       `db:"retry_times" json:"retry_times"` // 重试次数
	Status     int       `db:"status" json:"status"`           // 0:禁用 1:启用
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (ScheduledTask) TableName() string {
	return "sys_scheduled_tasks"
}

// TaskLog 任务执行日志表
type TaskLog struct {
	ID        uint      `db:"id" json:"id"`
	TaskID    uint      `db:"task_id" json:"task_id"`
	Status    int       `db:"status" json:"status"` // 0:失败 1:成功
	Result    string    `db:"result" json:"result"` // 执行结果
	Error     string    `db:"error" json:"error"`   // 错误信息
	StartTime time.Time `db:"start_time" json:"start_time"`
	EndTime   time.Time `db:"end_time" json:"end_time"`
}

func (TaskLog) TableName() string {
	return "sys_task_logs"
}

// ElementTrigger 元素触发器表
type ElementTrigger struct {
	ID           uint      `db:"id" json:"id"`
	AppID        uint      `db:"app_id" json:"app_id"`
	ElementType  string    `db:"element_type" json:"element_type"` // form:表单 model:数据模型
	ElementID    uint      `db:"element_id" json:"element_id"`
	TriggerPoint string    `db:"trigger_point" json:"trigger_point"` // before:之前 after:之后
	Type         string    `db:"type" json:"type"`                   // sql:SQL任务 http:HTTP任务
	Content      string    `db:"content" json:"content"`             // 任务内容
	Status       int       `db:"status" json:"status"`               // 0:禁用 1:启用
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (ElementTrigger) TableName() string {
	return "sys_element_triggers"
}
