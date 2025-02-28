package model

import (
	"database/sql"
)

// 字段信息
type Field struct {
	Name          string `json:"name" db:"name"`
	Comment       string `json:"comment" db:"comment"`
	ColumnType    string `json:"column_type" db:"column_type"`
	Sort          int    `json:"sort" db:"sort"`
	PrimaryKey    bool   `json:"primary_key,omitempty" db:"primary_key"`
	AutoIncrement bool   `json:"auto_increment,omitempty" db:"auto_increment"`
	NotNull       bool   `json:"not_null,omitempty" db:"not_null"`
	Default       string `json:"default,omitempty" db:"default"`
}

// FieldUpdateType 表示字段更新类型
type UpdateTypeString string

const (
	UpdateTypeAdd    UpdateTypeString = "add"
	UpdateTypeDrop   UpdateTypeString = "drop"
	UpdateTypeModify UpdateTypeString = "modify"
)

// FieldUpdate 表示字段更新信息
type FieldUpdate struct {
	UpdateType   UpdateTypeString // 更新类型：add, drop, modify
	OldFieldName string           // 旧字段名（用于修改字段时）
	Field        Field            // 新字段信息
}

type FieldUpdateReq struct {
	Fields []FieldUpdate `json:"fields"`
}

// IndexUpdate 表示索引更新信息
type IndexUpdate struct {
	UpdateType   UpdateTypeString // 更新类型：add, drop, modify
	OldIndexName string           // 旧索引名（用于修改索引时）
	Index        Index            // 新索引信息
}

type IndexUpdateReq struct {
	Indexes []IndexUpdate `json:"indexes"`
}

// 索引信息
type Index struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

// CreateTableReq 创建表请求
type CreateTableReq struct {
	ID          uint    `json:"id"`
	TableName   string  `json:"table_name"`
	DisplayName string  `json:"display_name"`
	Description string  `json:"description"`
	Func        string  `json:"func"`
	Fields      []Field `json:"fields"`
	Indexes     []Index `json:"indexes"`
	ParentID    uint    `json:"parent_id"`
}

// UpdateTableReq 更新表请求
type UpdateTableReq struct {
	// AppID       uint      `json:"app_id"`
	TableName   string `json:"table_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// TableUpdateReq 统一的表更新请求
type TableUpdateReq struct {
	TableName   string        `json:"table_name"`   // 表名
	DisplayName string        `json:"display_name"` // 显示名称
	Description string        `json:"description"`  // 描述
	Func        string        `json:"func"`         // 功能
	Fields      []FieldUpdate `json:"fields"`       // 字段更新
	Indexes     []IndexUpdate `json:"indexes"`      // 索引更新
}

// MySQLField 表示从 MySQL 获取的字段信息
// type MySQLField struct {
// 	Field      string `db:"Field"`
// 	Type       string `db:"Type"`
// 	Collation  string `db:"Collation"`
// 	Sort       int    `db:"Sort"`
// 	Null       string `db:"Null"`
// 	Key        string `db:"Key"`
// 	Default    string `db:"Default"`
// 	Extra      string `db:"Extra"`
// 	Privileges string `db:"Privileges"`
// 	Comment    string `db:"Comment"`
// }

// MySQLIndex 表示从 MySQL 获取的索引信息
type MySQLIndex struct {
	Table        string         `db:"Table"`
	NonUnique    int            `db:"Non_unique"`
	KeyName      string         `db:"Key_name"`
	SeqInIndex   int            `db:"Seq_in_index"`
	ColumnName   string         `db:"Column_name"`
	Collation    sql.NullString `db:"Collation"`
	Cardinality  sql.NullInt64  `db:"Cardinality"`
	SubPart      sql.NullInt64  `db:"Sub_part"`
	Packed       sql.NullString `db:"Packed"`
	Null         sql.NullString `db:"Null"`
	IndexType    string         `db:"Index_type"`
	Comment      sql.NullString `db:"Comment"`
	IndexComment sql.NullString `db:"Index_comment"`
	Visible      sql.NullString `db:"Visible"`
	Expression   sql.NullString `db:"Expression"`
}
