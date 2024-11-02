package model

import "time"

// User 用户模型
type User struct {
	TenantBase
	Username    string     `json:"username" db:"username"`
	Password    string     `json:"-" db:"password"` // 密码不返回给前端
	Salt        string     `json:"-" db:"salt"`     // 密码盐值
	Nickname    string     `json:"nickname" db:"nickname"`
	Email       string     `json:"email" db:"email"`
	Phone       string     `json:"phone" db:"phone"`
	Avatar      string     `json:"avatar" db:"avatar"`
	Status      int        `json:"status" db:"status"`               // 状态：0-禁用 1-启用
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"` // 最后登录时间
}

// UserRole 用户角色关联
type UserRole struct {
	TenantBase
	UserID uint `json:"user_id" db:"user_id"`
	RoleID uint `json:"role_id" db:"role_id"`
}

// UserLoginHistory 用户登录历史
type UserLoginHistory struct {
	Base
	UserID    uint      `json:"user_id" db:"user_id"`
	TenantID  uint      `json:"tenant_id" db:"tenant_id"`
	IP        string    `json:"ip" db:"ip"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	LoginAt   time.Time `json:"login_at" db:"login_at"`
	Status    int       `json:"status" db:"status"` // 状态：0-失败 1-成功
}

const (
	UserStatusDisabled = 0 // 禁用
	UserStatusEnabled  = 1 // 启用
)

const (
	LoginStatusFailed  = 0 // 登录失败
	LoginStatusSuccess = 1 // 登录成功
)
