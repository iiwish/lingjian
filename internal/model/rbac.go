package model

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	AppCode     string `json:"app_code" binding:"required"`
	ParentCode  string `json:"parent_code"`
	Description string `json:"description"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	AppCode     string `json:"app_code" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Description string `json:"description"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionCodes []string `json:"permission_codes" binding:"required"`
	AppCode         string   `json:"app_code" binding:"required"`
}
