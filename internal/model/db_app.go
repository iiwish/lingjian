package model

import "time"

// App 应用表
type App struct {
	ID          uint      `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Description string    `db:"description" json:"description"`
	Status      int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	CreatorID   uint      `db:"creator_id" json:"creator_id"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID   uint      `db:"updater_id" json:"updater_id"`
}

// SysVar 系统配置表
type SysVar struct {
	ID          uint      `db:"id" json:"id"`                   // 主键ID
	Name        string    `db:"name" json:"name"`               // 配置名称
	Code        string    `db:"code" json:"code"`               // 配置代码
	Value       string    `db:"value" json:"value"`             // 配置值
	Description string    `db:"description" json:"description"` // 配置描述
	Status      int       `db:"status" json:"status"`           // 状态 0:禁用 1:启用
	CreatedAt   time.Time `db:"created_at" json:"created_at"`   // 创建时间
	CreatorID   uint      `db:"creator_id" json:"creator_id"`   // 创建人ID
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`   // 更新时间
	UpdaterID   uint      `db:"updater_id" json:"updater_id"`   // 更新人ID
}
