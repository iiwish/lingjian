package config

import "github.com/iiwish/lingjian/internal/model"

// CreateTableRequest 创建数据表配置请求
type CreateTableRequest struct {
	AppID          uint               `json:"app_id" binding:"required"`
	Name           string             `json:"name" binding:"required"`
	Code           string             `json:"code" binding:"required"`
	Description    string             `json:"description"`
	MySQLTableName string             `json:"mysql_table_name" binding:"required"`
	Fields         []model.TableField `json:"fields" binding:"required"`
	Indexes        []model.TableIndex `json:"indexes"`
}

// CreateDimensionRequest 创建维度配置请求
type CreateDimensionRequest struct {
	AppID          uint                  `json:"app_id" binding:"required"`
	Name           string                `json:"name" binding:"required"`
	Code           string                `json:"code" binding:"required"`
	Type           string                `json:"type" binding:"required"`
	MySQLTableName string                `json:"mysql_table_name" binding:"required"`
	Configuration  model.DimensionConfig `json:"configuration" binding:"required"`
}

// CreateModelRequest 创建数据模型配置请求
type CreateModelRequest struct {
	AppID      uint                   `json:"app_id" binding:"required"`
	Name       string                 `json:"name" binding:"required"`
	Code       string                 `json:"code" binding:"required"`
	TableID    uint                   `json:"table_id" binding:"required"`
	Fields     []model.ModelField     `json:"fields" binding:"required"`
	Dimensions []model.ModelDimension `json:"dimensions"`
	Metrics    []model.ModelMetric    `json:"metrics"`
}

// CreateFormRequest 创建表单配置请求
type CreateFormRequest struct {
	AppID   uint              `json:"app_id" binding:"required"`
	Name    string            `json:"name" binding:"required"`
	Code    string            `json:"code" binding:"required"`
	Type    string            `json:"type" binding:"required"`
	TableID uint              `json:"table_id" binding:"required"`
	Layout  model.FormLayout  `json:"layout" binding:"required"`
	Fields  []model.FormField `json:"fields" binding:"required"`
	Rules   []model.FormRule  `json:"rules"`
	Events  []model.FormEvent `json:"events"`
}

// CreateMenuRequest 创建菜单配置请求
type CreateMenuRequest struct {
	AppID     uint   `json:"app_id" binding:"required"`
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Icon      string `json:"icon"`
	Path      string `json:"path" binding:"required"`
	Component string `json:"component" binding:"required"`
	Sort      int    `json:"sort"`
}
