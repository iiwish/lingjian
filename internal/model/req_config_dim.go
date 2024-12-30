package model

import (
	"github.com/iiwish/lingjian/pkg/utils"
)

// DimensionItem 维度配置
type DimensionItem struct {
	ID          uint             `db:"id" json:"id"`                   // 主键ID
	NodeID      string           `db:"node_id" json:"node_id"`         // 节点ID
	ParentID    uint             `db:"parent_id" json:"parent_id"`     // 父节点ID
	Name        string           `db:"name" json:"name"`               // 名称
	Code        string           `db:"code" json:"code"`               // 编码
	Description uint             `db:"description" json:"description"` // 维度ID
	Level       int              `db:"level" json:"level"`             // 层级
	Sort        int              `db:"sort" json:"sort"`               // 排序
	Status      int              `db:"status" json:"status"`           // 状态
	Custom1     string           `db:"custom1" json:"custom1"`         // 自定义字段1
	Custom2     string           `db:"custom2" json:"custom2"`         // 自定义字段2
	Custom3     string           `db:"custom3" json:"custom3"`         // 自定义字段3
	CreatedAt   utils.CustomTime `db:"created_at" json:"created_at"`   // 创建时间
	CreatorID   uint             `db:"creator_id" json:"creator_id"`   // 创建者ID
	UpdatedAt   utils.CustomTime `db:"updated_at" json:"updated_at"`   // 更新时间
	UpdaterID   uint             `db:"updater_id" json:"updater_id"`   // 更新者ID
}

// TreeDimensionItem 维度配置树形结构
type TreeDimensionItem struct {
	DimensionItem
	Children []*TreeDimensionItem `json:"children"` // 子维度列表
}

// CreateDimReq 创建维度请求
type CreateDimReq struct {
	TableName   string `db:"table_name" json:"table_name"`
	DisplayName string `db:"display_name" json:"display_name"`
	Description string `db:"description" json:"description"`
	ParentID    uint   `db:"parent_id" json:"parent_id"`
}
