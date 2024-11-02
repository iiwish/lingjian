package model

import "time"

// ConfigTable 数据表配置
type ConfigTable struct {
	ID             uint      `db:"id" json:"id"`
	ApplicationID  uint      `db:"application_id" json:"application_id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Description    string    `db:"description" json:"description"`
	MySQLTableName string    `db:"mysql_table_name" json:"mysql_table_name"` // 对应的MySQL表名
	Fields         string    `db:"fields" json:"fields"`                     // JSON格式，存储字段定义
	Indexes        string    `db:"indexes" json:"indexes"`                   // JSON格式，存储索引定义
	Status         int       `db:"status" json:"status"`                     // 0:禁用 1:启用
	Version        int       `db:"version" json:"version"`                   // 版本号
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// ConfigDimension 维度配置
type ConfigDimension struct {
	ID             uint      `db:"id" json:"id"`
	ApplicationID  uint      `db:"application_id" json:"application_id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Type           string    `db:"type" json:"type"`                         // time:时间维度 enum:枚举维度 range:范围维度
	MySQLTableName string    `db:"mysql_table_name" json:"mysql_table_name"` // 对应的MySQL表名
	Configuration  string    `db:"configuration" json:"configuration"`       // JSON格式，存储维度配置
	Status         int       `db:"status" json:"status"`                     // 0:禁用 1:启用
	Version        int       `db:"version" json:"version"`                   // 版本号
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// ConfigDataModel 数据模型配置
type ConfigDataModel struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	Name          string    `db:"name" json:"name"`
	Code          string    `db:"code" json:"code"`
	TableID       uint      `db:"table_id" json:"table_id"`     // 关联的数据表ID
	Fields        string    `db:"fields" json:"fields"`         // JSON格式，存储字段配置
	Dimensions    string    `db:"dimensions" json:"dimensions"` // JSON格式，存储维度配置
	Metrics       string    `db:"metrics" json:"metrics"`       // JSON格式，存储指标配置
	Status        int       `db:"status" json:"status"`         // 0:禁用 1:启用
	Version       int       `db:"version" json:"version"`       // 版本号
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ConfigForm 表单配置
type ConfigForm struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	Name          string    `db:"name" json:"name"`
	Code          string    `db:"code" json:"code"`
	Type          string    `db:"type" json:"type"`         // create:新建表单 edit:编辑表单 view:查看表单
	TableID       uint      `db:"table_id" json:"table_id"` // 关联的数据表ID
	Layout        string    `db:"layout" json:"layout"`     // JSON格式，存储表单布局
	Fields        string    `db:"fields" json:"fields"`     // JSON格式，存储字段配置
	Rules         string    `db:"rules" json:"rules"`       // JSON格式，存储验证规则
	Events        string    `db:"events" json:"events"`     // JSON格式，存储事件配置
	Status        int       `db:"status" json:"status"`     // 0:禁用 1:启用
	Version       int       `db:"version" json:"version"`   // 版本号
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ConfigMenu 菜单配置
type ConfigMenu struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	ParentID      uint      `db:"parent_id" json:"parent_id"` // 父菜单ID，0表示顶级菜单
	Name          string    `db:"name" json:"name"`
	Code          string    `db:"code" json:"code"`
	Icon          string    `db:"icon" json:"icon"`           // 菜单图标
	Path          string    `db:"path" json:"path"`           // 菜单路径
	Component     string    `db:"component" json:"component"` // 关联的前端组件
	Sort          int       `db:"sort" json:"sort"`           // 排序号
	Status        int       `db:"status" json:"status"`       // 0:禁用 1:启用
	Version       int       `db:"version" json:"version"`     // 版本号
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ConfigVersion 配置版本记录
type ConfigVersion struct {
	ID            uint      `db:"id" json:"id"`
	ApplicationID uint      `db:"application_id" json:"application_id"`
	ConfigType    string    `db:"config_type" json:"config_type"` // table:数据表 dimension:维度 model:数据模型 form:表单 menu:菜单
	ConfigID      uint      `db:"config_id" json:"config_id"`     // 配置ID
	Version       int       `db:"version" json:"version"`         // 版本号
	Content       string    `db:"content" json:"content"`         // JSON格式，存储配置内容
	Comment       string    `db:"comment" json:"comment"`         // 版本说明
	CreatorID     uint      `db:"creator_id" json:"creator_id"`   // 创建人ID
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
