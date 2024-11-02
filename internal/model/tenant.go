package model

// Tenant 租户模型
type Tenant struct {
	Base
	Name        string `json:"name" db:"name"`               // 租户名称
	Code        string `json:"code" db:"code"`               // 租户编码
	Description string `json:"description" db:"description"` // 租户描述
	Status      int    `json:"status" db:"status"`           // 状态：0-禁用 1-启用
	AdminEmail  string `json:"admin_email" db:"admin_email"` // 管理员邮箱
	AdminPhone  string `json:"admin_phone" db:"admin_phone"` // 管理员电话
	Logo        string `json:"logo" db:"logo"`               // 租户logo
	Domain      string `json:"domain" db:"domain"`           // 租户域名
	Config      string `json:"config" db:"config"`           // 租户配置（JSON）
}

// TenantDatabase 租户数据库配置
type TenantDatabase struct {
	Base
	TenantID uint   `json:"tenant_id" db:"tenant_id"`
	Host     string `json:"host" db:"host"`         // 数据库主机
	Port     int    `json:"port" db:"port"`         // 数据库端口
	Database string `json:"database" db:"database"` // 数据库名
	Username string `json:"username" db:"username"` // 用户名
	Password string `json:"-" db:"password"`        // 密码
	Status   int    `json:"status" db:"status"`     // 状态：0-禁用 1-启用
}

// TenantConfig 租户配置
type TenantConfig struct {
	Theme    string            `json:"theme"`    // 主题
	Modules  []string          `json:"modules"`  // 启用的模块
	Settings map[string]string `json:"settings"` // 其他设置
}

const (
	TenantStatusDisabled = 0 // 禁用
	TenantStatusEnabled  = 1 // 启用
)

const (
	DBStatusDisabled = 0 // 禁用
	DBStatusEnabled  = 1 // 启用
)
