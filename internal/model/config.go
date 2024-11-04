package model

import "time"

// ConfigTable 数据表配置
type ConfigTable struct {
	ID             uint      `db:"id" json:"id"`
	ApplicationID  uint      `db:"application_id" json:"application_id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Description    string    `db:"description" json:"description"`
	MySQLTableName string    `db:"mysql_table_name" json:"mysql_table_name"`
	Fields         string    `db:"fields" json:"fields"`   // JSON格式，存储字段定义
	Indexes        string    `db:"indexes" json:"indexes"` // JSON格式，存储索引定义
	Status         int       `db:"status" json:"status"`   // 0:禁用 1:启用
	Version        int       `db:"version" json:"version"` // 版本号
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// TableField 表字段定义
type TableField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Length      int    `json:"length,omitempty"`
	Nullable    bool   `json:"nullable"`
	Default     any    `json:"default,omitempty"`
	Comment     string `json:"comment,omitempty"`
	IsPrimary   bool   `json:"is_primary,omitempty"`
	IsAutoInc   bool   `json:"is_auto_inc,omitempty"`
	IsUnique    bool   `json:"is_unique,omitempty"`
	UniqueGroup string `json:"unique_group,omitempty"`
}

// TableIndex 表索引定义
type TableIndex struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"` // NORMAL, UNIQUE, FULLTEXT
	Columns []string `json:"columns"`
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

// DimensionConfig 维度配置详情
type DimensionConfig struct {
	NodeIDField   string `json:"node_id_field"`   // 节点ID字段
	ParentIDField string `json:"parent_id_field"` // 父节点ID字段
	NameField     string `json:"name_field"`      // 名称字段
	CodeField     string `json:"code_field"`      // 编码字段
	LevelField    string `json:"level_field"`     // 层级字段
	OrderField    string `json:"order_field"`     // 排序字段
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

// ModelField 模型字段配置
type ModelField struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	TableField  string `json:"table_field"`
}

// ModelDimension 模型维度配置
type ModelDimension struct {
	DimensionID  uint   `json:"dimension_id"`
	JoinField    string `json:"join_field"`
	DisplayField string `json:"display_field"`
}

// ModelMetric 模型指标配置
type ModelMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // sum, avg, count, etc.
	Expression  string `json:"expression"`
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

// FormLayout 表单布局配置
type FormLayout struct {
	Type    string       `json:"type"` // grid, flex等
	Columns []FormColumn `json:"columns"`
}

// FormColumn 表单列配置
type FormColumn struct {
	Span     int      `json:"span"`
	Elements []string `json:"elements"` // 引用FormField的ID
}

// FormField 表单字段配置
type FormField struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"` // text, input, select等
	Label      string         `json:"label"`
	Field      string         `json:"field"`
	Required   bool           `json:"required"`
	Properties map[string]any `json:"properties"`
}

// FormRule 表单验证规则
type FormRule struct {
	Field    string `json:"field"`
	Type     string `json:"type"` // required, email, mobile等
	Message  string `json:"message"`
	Pattern  string `json:"pattern"`  // 正则表达式
	Min      *int   `json:"min"`      // 最小值/长度
	Max      *int   `json:"max"`      // 最大值/长度
	Trigger  string `json:"trigger"`  // blur, change等
	Required bool   `json:"required"` // 是否必填
}

// FormEvent 表单事件配置
type FormEvent struct {
	Type       string         `json:"type"`       // before_submit, after_submit等
	Action     string         `json:"action"`     // sql, api等
	Content    string         `json:"content"`    // SQL语句或API地址
	Parameters map[string]any `json:"parameters"` // 参数配置
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
