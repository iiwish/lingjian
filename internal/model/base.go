package model

import "time"

// Base 基础模型
type Base struct {
	ID        uint       `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// TenantBase 带租户的基础模型
type TenantBase struct {
	Base
	TenantID uint `json:"tenant_id" db:"tenant_id"`
}
