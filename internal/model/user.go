package model

import "time"

// User 用户表
type User struct {
	ID        uint      `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Nickname  string    `db:"nickname" json:"nickname"`
	Avatar    string    `db:"avatar" json:"avatar"`
	Password  string    `db:"password" json:"-"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	Status    int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatorID uint      `db:"creator_id" json:"creator_id"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID uint      `db:"updater_id" json:"updater_id"`
}

// Role 角色表
type Role struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	Status    int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatorID uint      `db:"creator_id" json:"creator_id"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID uint      `db:"updater_id" json:"updater_id"`
}

// Permission 权限表
type Permission struct {
	ID          uint      `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Type        string    `db:"type" json:"type"`
	Path        string    `db:"path" json:"path"`
	Method      string    `db:"method" json:"method"`
	Status      int       `db:"status" json:"status"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	CreatorID   uint      `db:"creator_id" json:"creator_id"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID   uint      `db:"updater_id" json:"updater_id"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint      `db:"id" json:"id"`
	UserID    uint      `db:"user_id" json:"user_id"`
	RoleID    uint      `db:"role_id" json:"role_id"`
	CreateAt  time.Time `db:"created_at" json:"created_at"`
	CreatorID uint      `db:"creator_id" json:"creator_id"`
	UpdateAt  time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID uint      `db:"updater_id" json:"updater_id"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `db:"id" json:"id"`
	RoleID       uint      `db:"role_id" json:"role_id"`
	PermissionID uint      `db:"permission_id" json:"permission_id"`
	CreateAt     time.Time `db:"created_at" json:"created_at"`
	CreatorID    uint      `db:"creator_id" json:"creator_id"`
	UpdateAt     time.Time `db:"updated_at" json:"updated_at"`
	UpdaterID    uint      `db:"updater_id" json:"updater_id"`
}
