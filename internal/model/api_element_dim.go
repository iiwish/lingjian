package model

import "github.com/iiwish/lingjian/pkg/utils"

// DimensionItem 维度数据项
type DimensionItem struct {
	ID          uint              `db:"id" json:"id"`                   // 主键ID
	NodeID      string            `db:"node_id" json:"node_id"`         // 节点ID
	ParentID    uint              `db:"parent_id" json:"parent_id"`     // 父节点ID
	Name        string            `db:"name" json:"name"`               // 名称
	Code        string            `db:"code" json:"code"`               // 编码
	Description string            `db:"description" json:"description"` // 描述
	Level       int               `db:"level" json:"level"`             // 层级
	Sort        int               `db:"sort" json:"sort"`               // 排序
	Status      int               `db:"status" json:"status"`           // 状态
	CustomData  map[string]string `json:"custom_data"`                  // 自定义列数据
	CreatedAt   utils.CustomTime  `db:"created_at" json:"created_at"`   // 创建时间
	CreatorID   uint              `db:"creator_id" json:"creator_id"`   // 创建者ID
	UpdatedAt   utils.CustomTime  `db:"updated_at" json:"updated_at"`   // 更新时间
	UpdaterID   uint              `db:"updater_id" json:"updater_id"`   // 更新者ID
}

// TreeDimensionItem 维度数据树形结构
type TreeDimensionItem struct {
	DimensionItem
	Children []*TreeDimensionItem `json:"children"` // 子维度列表
}

// CreateDimensionItemReq 创建维度数据请求
type CreateDimensionItemReq struct {
	ParentID    uint              `json:"parent_id"`   // 父节点ID
	Name        string            `json:"name"`        // 名称
	Code        string            `json:"code"`        // 编码
	Description string            `json:"description"` // 描述
	CustomData  map[string]string `json:"custom_data"` // 自定义列数据
}

// UpdateDimensionItemReq 更新维度数据请求
type UpdateDimensionItemReq struct {
	ID          uint              `json:"id"`          // 主键ID
	Name        string            `json:"name"`        // 名称
	Code        string            `json:"code"`        // 编码
	Description string            `json:"description"` // 描述
	CustomData  map[string]string `json:"custom_data"` // 自定义列数据
}

// BatchCreateDimensionItemReq 批量创建维度数据请求
type BatchCreateDimensionItemReq struct {
	Items []CreateDimensionItemReq `json:"items"` // 维度数据项列表
}

// BatchUpdateDimensionItemReq 批量更新维度数据请求
type BatchUpdateDimensionItemReq struct {
	Items []UpdateDimensionItemReq `json:"items"` // 维度数据项列表
}
