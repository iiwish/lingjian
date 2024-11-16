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
