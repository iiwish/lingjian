package model

import (
	"github.com/iiwish/lingjian/pkg/utils"
)

// User 用户表
type User struct {
	ID        uint             `db:"id" json:"id"`
	Username  string           `db:"username" json:"username"`
	Nickname  string           `db:"nickname" json:"nickname"`
	Avatar    string           `db:"avatar" json:"avatar"`
	Password  string           `db:"password" json:"password"`
	Email     string           `db:"email" json:"email"`
	Phone     string           `db:"phone" json:"phone"`
	Status    int              `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID uint             `db:"updater_id" json:"updater_id"`
	DeletedAt utils.CustomTime `db:"deleted_at" json:"deleted_at"`
}

// Role 角色表
type Role struct {
	ID          uint             `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Code        string           `db:"code" json:"code"`
	ParentID    uint             `db:"parent_id" json:"parent_id"`
	Description string           `db:"description" json:"description"`
	Status      int              `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt   utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID   uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt   utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID   uint             `db:"updater_id" json:"updater_id"`
	DeletedAt   utils.CustomTime `db:"deleted_at" json:"deleted_at"`
}

// Permission 权限表
type Permission struct {
	ID          uint             `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Code        string           `db:"code" json:"code"`
	Type        string           `db:"type" json:"type"`
	Path        string           `db:"path" json:"path"`
	Method      string           `db:"method" json:"method"`
	DimID       uint             `db:"dim_id" json:"dim_id"`
	ItemID      uint             `db:"item_id" json:"item_id"`
	Status      int              `db:"status" json:"status"`
	Description string           `db:"description" json:"description"`
	CreatedAt   utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID   uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt   utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID   uint             `db:"updater_id" json:"updater_id"`
	DeletedAt   utils.CustomTime `db:"deleted_at" json:"deleted_at"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint             `db:"id" json:"id"`
	UserID    uint             `db:"user_id" json:"user_id"`
	RoleID    uint             `db:"role_id" json:"role_id"`
	CreateAt  utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID uint             `db:"creator_id" json:"creator_id"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint             `db:"id" json:"id"`
	RoleID       uint             `db:"role_id" json:"role_id"`
	PermissionID uint             `db:"permission_id" json:"permission_id"`
	CreateAt     utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID    uint             `db:"creator_id" json:"creator_id"`
}
