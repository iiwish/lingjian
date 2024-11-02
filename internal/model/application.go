package model

// Application 应用模型
type Application struct {
	TenantBase
	Name        string `json:"name" db:"name"`               // 应用名称
	Code        string `json:"code" db:"code"`               // 应用编码
	Description string `json:"description" db:"description"` // 应用描述
	Status      int    `json:"status" db:"status"`           // 状态：0-禁用 1-启用
	Type        string `json:"type" db:"type"`               // 应用类型
	Icon        string `json:"icon" db:"icon"`               // 应用图标
	Config      string `json:"config" db:"config"`           // 应用配置（JSON）
	IsTemplate  bool   `json:"is_template" db:"is_template"` // 是否为模板
}

// ApplicationConfig 应用配置
type ApplicationConfig struct {
	Theme    string                 `json:"theme"`    // 主题
	Layout   string                 `json:"layout"`   // 布局
	Menu     []ApplicationMenu      `json:"menu"`     // 菜单配置
	Database []ApplicationDatabase  `json:"database"` // 数据库配置
	Settings map[string]interface{} `json:"settings"` // 其他设置
}

// ApplicationMenu 应用菜单
type ApplicationMenu struct {
	Name       string            `json:"name"`       // 菜单名称
	Path       string            `json:"path"`       // 路由路径
	Component  string            `json:"component"`  // 组件路径
	Icon       string            `json:"icon"`       // 图标
	Sort       int               `json:"sort"`       // 排序
	ParentID   *uint             `json:"parent_id"`  // 父级ID
	Permission string            `json:"permission"` // 权限编码
	Meta       map[string]string `json:"meta"`       // 元数据
}

// ApplicationDatabase 应用数据库表配置
type ApplicationDatabase struct {
	TableName   string             `json:"table_name"`   // 表名
	Description string             `json:"description"`  // 描述
	Fields      []ApplicationField `json:"fields"`       // 字段配置
	Indexes     []ApplicationIndex `json:"indexes"`      // 索引配置
	TablePrefix string             `json:"table_prefix"` // 表前缀
}

// ApplicationField 应用数据库字段配置
type ApplicationField struct {
	Name       string                 `json:"name"`       // 字段名
	Type       string                 `json:"type"`       // 字段类型
	Length     int                    `json:"length"`     // 字段长度
	Required   bool                   `json:"required"`   // 是否必填
	Unique     bool                   `json:"unique"`     // 是否唯一
	Default    interface{}            `json:"default"`    // 默认值
	Comment    string                 `json:"comment"`    // 注释
	Validation map[string]interface{} `json:"validation"` // 验证规则
}

// ApplicationIndex 应用数据库索引配置
type ApplicationIndex struct {
	Name   string   `json:"name"`   // 索引名
	Fields []string `json:"fields"` // 索引字段
	Type   string   `json:"type"`   // 索引类型：normal, unique
}

const (
	AppStatusDisabled = 0 // 禁用
	AppStatusEnabled  = 1 // 启用
)

const (
	AppTypeNormal   = "normal"   // 普通应用
	AppTypeTemplate = "template" // 模板应用
)

const (
	IndexTypeNormal = "normal" // 普通索引
	IndexTypeUnique = "unique" // 唯一索引
)
