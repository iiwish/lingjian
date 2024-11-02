package model

// Role 角色模型
type Role struct {
	TenantBase
	Name        string `json:"name" db:"name"`               // 角色名称
	Code        string `json:"code" db:"code"`               // 角色编码
	Description string `json:"description" db:"description"` // 角色描述
	Status      int    `json:"status" db:"status"`           // 状态：0-禁用 1-启用
	IsSystem    bool   `json:"is_system" db:"is_system"`     // 是否系统角色
}

// Permission 权限模型
type Permission struct {
	TenantBase
	Name        string `json:"name" db:"name"`               // 权限名称
	Code        string `json:"code" db:"code"`               // 权限编码
	Type        string `json:"type" db:"type"`               // 权限类型：menu-菜单 button-按钮 api-接口
	ParentID    *uint  `json:"parent_id" db:"parent_id"`     // 父级ID
	Path        string `json:"path" db:"path"`               // 路径
	Method      string `json:"method" db:"method"`           // HTTP方法
	Component   string `json:"component" db:"component"`     // 前端组件
	Icon        string `json:"icon" db:"icon"`               // 图标
	Sort        int    `json:"sort" db:"sort"`               // 排序
	Status      int    `json:"status" db:"status"`           // 状态：0-禁用 1-启用
	IsSystem    bool   `json:"is_system" db:"is_system"`     // 是否系统权限
	Description string `json:"description" db:"description"` // 权限描述
}

// RolePermission 角色权限关联
type RolePermission struct {
	TenantBase
	RoleID       uint `json:"role_id" db:"role_id"`
	PermissionID uint `json:"permission_id" db:"permission_id"`
}

const (
	PermTypeMenu   = "menu"   // 菜单权限
	PermTypeButton = "button" // 按钮权限
	PermTypeAPI    = "api"    // 接口权限
)

const (
	RoleStatusDisabled = 0 // 禁用
	RoleStatusEnabled  = 1 // 启用
)

const (
	PermStatusDisabled = 0 // 禁用
	PermStatusEnabled  = 1 // 启用
)
