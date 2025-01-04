package model

import "github.com/iiwish/lingjian/pkg/utils"

// TreeMenuItem 返回给前端的菜单树形结构
type TreeMenuItem struct {
	ID        uint             `db:"id" json:"id"`
	NodeID    string           `db:"node_id" json:"node_id"`
	ParentID  uint             `db:"parent_id" json:"parent_id"`
	MenuName  string           `db:"menu_name" json:"menu_name"`
	MenuCode  string           `db:"menu_code" json:"menu_code"`
	MenuType  int              `db:"menu_type" json:"menu_type"`
	Level     int              `db:"level" json:"level"`
	Sort      int              `db:"sort" json:"sort"`
	IconPath  string           `db:"icon_path" json:"icon_path"`
	SourceID  uint             `db:"source_id" json:"source_id"`
	Status    int              `db:"status" json:"status"`
	CreatedAt utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID uint             `db:"updater_id" json:"updater_id"`
	Children  []*TreeMenuItem  `json:"children"` // 子菜单列表
}

type CreateMenuItemReq struct {
	MenuName    string `db:"menu_name" json:"menu_name"`     // 名称
	MenuCode    string `db:"menu_code" json:"menu_code"`     // 编码
	Description string `db:"description" json:"description"` // 描述
	Status      int    `db:"status" json:"status"`           // 状态
	SourceID    uint   `db:"source_id" json:"source_id"`     // 菜单图标
	MenuType    int    `db:"menu_type" json:"menu_type"`     // 菜单类型
	IconPath    string `db:"icon_path" json:"icon_path"`     // 菜单图标
	ParentID    uint   `db:"parent_id" json:"parent_id"`     // 父节点ID
}

type UpdateMenuItemReq struct {
	MenuName    string `db:"menu_name" json:"menu_name"`     // 名称
	MenuCode    string `db:"menu_code" json:"menu_code"`     // 编码
	Description string `db:"description" json:"description"` // 描述
	Status      int    `db:"status" json:"status"`           // 状态
	SourceID    uint   `db:"source_id" json:"source_id"`     // 菜单图标
	MenuType    int    `db:"menu_type" json:"menu_type"`     // 菜单类型
	IconPath    string `db:"icon_path" json:"icon_path"`     // 菜单图标
	ID          uint   `db:"id" json:"id"`                   // 主键ID
}
