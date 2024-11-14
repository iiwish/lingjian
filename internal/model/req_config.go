package model

import "time"

type TreeConfigMenu struct {
	ID        uint              `db:"id" json:"id"`
	AppID     uint              `db:"app_id" json:"app_id"`
	NodeID    string            `db:"node_id" json:"node_id"`
	ParentID  uint              `db:"parent_id" json:"parent_id"`
	MenuName  string            `db:"menu_name" json:"menu_name"`
	MenuCode  string            `db:"menu_code" json:"menu_code"`
	MenuType  string            `db:"menu_type" json:"menu_type"`
	Level     int               `db:"level" json:"level"`
	Sort      int               `db:"sort" json:"sort"`
	Icon      string            `db:"icon" json:"icon"`
	Path      string            `db:"path" json:"path"`
	Status    int               `db:"status" json:"status"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	CreatorID uint              `db:"creator_id" json:"creator_id"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
	UpdaterID uint              `db:"updater_id" json:"updater_id"`
	Children  []*TreeConfigMenu `json:"children"` // 子菜单列表
}
