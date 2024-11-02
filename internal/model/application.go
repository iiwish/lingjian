package model

import "time"

// Application 应用表
type Application struct {
	ID          uint      `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Description string    `db:"description" json:"description"`
	Status      int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// UserApplication 用户应用关联表
type UserApplication struct {
	ID            uint      `db:"id" json:"id"`
	UserID        uint      `db:"user_id" json:"user_id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	IsDefault     bool      `db:"is_default" json:"is_default"` // 是否默认应用
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ApplicationTemplate 应用模板表
type ApplicationTemplate struct {
	ID            uint      `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	Configuration string    `db:"configuration" json:"configuration"` // JSON格式的配置信息
	Price         float64   `db:"price" json:"price"`                 // 模板价格，0表示免费
	CreatorID     uint      `db:"creator_id" json:"creator_id"`       // 创建者ID
	Downloads     int       `db:"downloads" json:"downloads"`         // 下载次数
	Status        int       `db:"status" json:"status"`               // 0:未上架 1:已上架
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
