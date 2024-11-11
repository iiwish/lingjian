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
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "sys_users"
}

// Role 角色表
type Role struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	Status    int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Role) TableName() string {
	return "sys_roles"
}

// Permission 权限表
type Permission struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	Type      string    `db:"type" json:"type"` // menu:菜单 api:接口
	Path      string    `db:"path" json:"path"`
	Method    string    `db:"method" json:"method"`
	Status    int       `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Permission) TableName() string {
	return "sys_permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `db:"id" json:"id"`
	UserID uint `db:"user_id" json:"user_id"`
	RoleID uint `db:"role_id" json:"role_id"`
}

func (UserRole) TableName() string {
	return "sys_user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint `db:"id" json:"id"`
	RoleID       uint `db:"role_id" json:"role_id"`
	PermissionID uint `db:"permission_id" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "sys_role_permissions"
}
