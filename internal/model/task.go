package model

import "time"

// ScheduledTask 定时任务表
type ScheduledTask struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	Name          string    `db:"name" json:"name"`
	Type          string    `db:"type" json:"type"`               // sql:SQL任务 http:HTTP任务
	Cron          string    `db:"cron" json:"cron"`               // cron表达式
	Content       string    `db:"content" json:"content"`         // 任务内容（SQL语句或HTTP配置）
	Timeout       int       `db:"timeout" json:"timeout"`         // 超时时间（秒）
	RetryTimes    int       `db:"retry_times" json:"retry_times"` // 重试次数
	Status        int       `db:"status" json:"status"`           // 0:禁用 1:启用
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
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

// ElementTrigger 元素触发器表
type ElementTrigger struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	ElementType   string    `db:"element_type" json:"element_type"` // form:表单 model:数据模型
	ElementID     uint      `db:"element_id" json:"element_id"`
	TriggerPoint  string    `db:"trigger_point" json:"trigger_point"` // before:之前 after:之后
	Type          string    `db:"type" json:"type"`                   // sql:SQL任务 http:HTTP任务
	Content       string    `db:"content" json:"content"`             // 任务内容
	Status        int       `db:"status" json:"status"`               // 0:禁用 1:启用
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
