package model

import "time"

// ConfigMenu 返回给前端的菜单树形结构
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

// 字段信息
type Field struct {
	Name          string `json:"name"`
	Comment       string `json:"comment"`
	Type          string `json:"type"`
	Sort          int    `json:"sort"`
	PrimaryKey    bool   `json:"primary_key,omitempty"`
	AutoIncrement bool   `json:"auto_increment,omitempty"`
	NotNull       bool   `json:"not_null,omitempty"`
	Default       string `json:"default,omitempty"`
}

// 索引信息
type Index struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

// CreateTableReq 创建表请求
type CreateTableReq struct {
	AppID       uint    `json:"app_id"`
	TableName   string  `json:"table_name"`
	DisplayName string  `json:"display_name"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
	Indexes     []Index `json:"indexes"`
}

// UpdateTableReq 更新表请求
type UpdateTableReq struct {
	// AppID       uint      `json:"app_id"`
	TableName   string `json:"table_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// MySQLField 表示从 MySQL 获取的字段信息
type MySQLField struct {
	Field      string `db:"Field"`
	Type       string `db:"Type"`
	Collation  string `db:"Collation"`
	Null       string `db:"Null"`
	Key        string `db:"Key"`
	Default    string `db:"Default"`
	Extra      string `db:"Extra"`
	Privileges string `db:"Privileges"`
	Comment    string `db:"Comment"`
}

// MySQLIndex 表示从 MySQL 获取的索引信息
type MySQLIndex struct {
	Table        string `db:"Table"`
	NonUnique    int    `db:"Non_unique"`
	KeyName      string `db:"Key_name"`
	SeqInIndex   int    `db:"Seq_in_index"`
	ColumnName   string `db:"Column_name"`
	Collation    string `db:"Collation"`
	Cardinality  int    `db:"Cardinality"`
	SubPart      int    `db:"Sub_part"`
	Packed       string `db:"Packed"`
	Null         string `db:"Null"`
	IndexType    string `db:"Index_type"`
	Comment      string `db:"Comment"`
	IndexComment string `db:"Index_comment"`
}
