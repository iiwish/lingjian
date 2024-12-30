package model

import "github.com/iiwish/lingjian/pkg/utils"

// ConfigTable 数据表配置
type ConfigTable struct {
	ID          uint             `db:"id" json:"id"`
	AppID       uint             `db:"app_id" json:"app_id"`
	TableName   string           `db:"table_name" json:"table_name"`
	DisplayName string           `db:"display_name" json:"display_name"`
	Description string           `db:"description" json:"description"`
	Func        string           `db:"func" json:"func"`
	Status      int              `db:"status" json:"status"` // 0:禁用 1:启用
	CreatedAt   utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID   uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt   utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID   uint             `db:"updater_id" json:"updater_id"`
}

// ConfigDimension 维度配置
type ConfigDimension struct {
	ID            uint             `db:"id" json:"id"`
	AppID         uint             `db:"app_id" json:"app_id"`
	TableName     string           `db:"table_name" json:"table_name"`
	DisplayName   string           `db:"display_name" json:"display_name"`
	Description   string           `db:"description" json:"description"`
	Status        int              `db:"status" json:"status"` // 0:禁用 1:启用
	CustomColumns string           `db:"custom_columns" json:"custom_columns"`
	CreatedAt     utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID     uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt     utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID     uint             `db:"updater_id" json:"updater_id"`
}

// ConfigModel 数据模型配置
type ConfigModel struct {
	ID            uint             `db:"id" json:"id"`
	AppID         uint             `db:"app_id" json:"app_id"`
	ModelName     string           `db:"model_name" json:"model_name"`
	DisplayName   string           `db:"display_name" json:"display_name"`
	Description   string           `db:"description" json:"description"`
	Configuration string           `db:"configuration" json:"configuration"`
	Status        int              `db:"status" json:"status"`   // 0:禁用 1:启用
	Version       int              `db:"version" json:"version"` // 版本号
	CreatedAt     utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID     uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt     utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID     uint             `db:"updater_id" json:"updater_id"`
}

// ConfigForm 表单配置
type ConfigForm struct {
	ID            uint             `db:"id" json:"id"`
	AppID         uint             `db:"app_id" json:"app_id"`
	ModelID       uint             `db:"model_id" json:"model_id"`
	FormName      string           `db:"form_name" json:"form_name"`
	FormType      string           `db:"form_type" json:"form_type"`
	DisplayName   string           `db:"display_name" json:"display_name"`
	Description   string           `db:"description" json:"description"`
	Configuration string           `db:"configuration" json:"configuration"`
	Status        int              `db:"status" json:"status"`   // 0:禁用 1:启用
	Version       int              `db:"version" json:"version"` // 版本号
	CreatedAt     utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID     uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt     utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID     uint             `db:"updater_id" json:"updater_id"`
}

// ConfigMenu 菜单配置
type ConfigMenu struct {
	ID        uint             `db:"id" json:"id"`
	AppID     uint             `db:"app_id" json:"app_id"`
	NodeID    string           `db:"node_id" json:"node_id"`     // 节点ID
	ParentID  uint             `db:"parent_id" json:"parent_id"` // 父菜单ID，0表示顶级菜单
	MenuName  string           `db:"menu_name" json:"menu_name"` // 菜单名称
	MenuCode  string           `db:"menu_code" json:"menu_code"` // 菜单代码
	MenuType  int              `db:"menu_type" json:"menu_type"` // 菜单类型
	Level     int              `db:"level" json:"level"`         // 菜单层级
	Sort      int              `db:"sort" json:"sort"`           // 排序号
	Icon      string           `db:"icon" json:"icon"`           // 菜单图标
	SourceID  uint             `db:"source_id" json:"source_id"` // 菜单图标
	Status    int              `db:"status" json:"status"`       // 0:禁用 1:启用
	CreatedAt utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID uint             `db:"updater_id" json:"updater_id"`
}

// ConfigVersion 配置版本记录
type ConfigVersion struct {
	ID         uint             `db:"id" json:"id"`
	AppID      uint             `db:"app_id" json:"app_id"`
	ConfigType string           `db:"config_type" json:"config_type"` // table:数据表 dimension:维度 model:数据模型 form:表单
	ConfigID   uint             `db:"config_id" json:"config_id"`     // 配置ID
	Version    int              `db:"version" json:"version"`         // 版本号
	Comment    string           `db:"comment" json:"comment"`         // 版本说明
	CreatedAt  utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID  uint             `db:"creator_id" json:"creator_id"` // 创建人ID
	UpdateAt   utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID  uint             `db:"updater_id" json:"updater_id"` // 更新人ID
}
