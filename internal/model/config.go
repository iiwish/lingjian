package model

import "time"

// ConfigTable 数据表配置
type ConfigTable struct {
	ID          uint64     `db:"id" json:"id"`
	AppID       uint64     `db:"app_id" json:"app_id"`
	TableName   string     `db:"table_name" json:"table_name"`
	DisplayName string     `db:"display_name" json:"display_name"`
	Description string     `db:"description" json:"description"`
	Filters     string     `db:"filters" json:"filters"` // JSON string
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

// TableFilter 表格筛选配置
type TableFilter struct {
	ColumnName string `json:"column_name"`
	FilterType string `json:"filter_type"` // time或other
}

// ConfigDimension 维度配置
type ConfigDimension struct {
	ID          uint64     `db:"id" json:"id"`
	AppID       uint64     `db:"app_id" json:"app_id"`
	TableName   string     `db:"table_name" json:"table_name"`
	DisplayName string     `db:"display_name" json:"display_name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

// ConfigModel 数据模型配置
type ConfigModel struct {
	ID            uint64     `db:"id" json:"id"`
	AppID         uint64     `db:"app_id" json:"app_id"`
	ModelName     string     `db:"model_name" json:"model_name"`
	DisplayName   string     `db:"display_name" json:"display_name"`
	Description   string     `db:"description" json:"description"`
	Configuration string     `db:"configuration" json:"configuration"` // JSON string
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

// ModelConfiguration 数据模型配置详情
type ModelConfiguration struct {
	Tables []ModelTable `json:"tables"`
}

// ModelTable 数据模型中的表配置
type ModelTable struct {
	TableName string          `json:"table_name"`
	Relations []TableRelation `json:"relations"`
}

// TableRelation 表关系配置
type TableRelation struct {
	TargetTable  string `json:"target_table"`
	SourceColumn string `json:"source_column"`
	TargetColumn string `json:"target_column"`
	RelationType string `json:"relation_type"` // one_to_one, one_to_many
}

// ConfigForm 表单配置
type ConfigForm struct {
	ID            uint64     `db:"id" json:"id"`
	AppID         uint64     `db:"app_id" json:"app_id"`
	ModelID       uint64     `db:"model_id" json:"model_id"`
	FormName      string     `db:"form_name" json:"form_name"`
	DisplayName   string     `db:"display_name" json:"display_name"`
	Description   string     `db:"description" json:"description"`
	Configuration string     `db:"configuration" json:"configuration"` // JSON string
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

// FormConfiguration 表单配置详情
type FormConfiguration struct {
	Elements []FormElement `json:"elements"`
	Layout   FormLayout    `json:"layout"`
}

// FormElement 表单元素
type FormElement struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"` // text, input, table等
	Label      string      `json:"label"`
	Field      string      `json:"field"`
	Properties interface{} `json:"properties"`
}

// FormLayout 表单布局
type FormLayout struct {
	Type    string       `json:"type"` // grid, flex等
	Columns []FormColumn `json:"columns"`
}

// FormColumn 表单列
type FormColumn struct {
	Span     int      `json:"span"`
	Elements []string `json:"elements"` // 引用FormElement的ID
}

// ConfigMenu 菜单配置
type ConfigMenu struct {
	ID         uint64     `db:"id" json:"id"`
	AppID      uint64     `db:"app_id" json:"app_id"`
	NodeID     string     `db:"node_id" json:"node_id"`
	ParentID   string     `db:"parent_id" json:"parent_id"`
	MenuName   string     `db:"menu_name" json:"menu_name"`
	MenuCode   string     `db:"menu_code" json:"menu_code"`
	MenuType   int8       `db:"menu_type" json:"menu_type"`
	Level      int        `db:"level" json:"level"`
	Sort       int        `db:"sort" json:"sort"`
	Icon       string     `db:"icon" json:"icon"`
	Path       string     `db:"path" json:"path"`
	Component  string     `db:"component" json:"component"`
	Permission string     `db:"permission" json:"permission"`
	Visible    bool       `db:"visible" json:"visible"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}
