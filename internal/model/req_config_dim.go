package model

import (
	"fmt"

	"github.com/iiwish/lingjian/pkg/utils"
)

// CustomColumn 自定义列定义
type CustomColumn struct {
	Name    string `json:"name"`    // 列名
	Comment string `json:"comment"` // 列注释
	Length  int    `json:"length"`  // 长度
}

// Validate 验证自定义列定义
func (c *CustomColumn) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	// 验证列名格式(只允许字母数字下划线,且以字母开头)
	if !utils.IsValidIdentifier(c.Name) {
		return fmt.Errorf("invalid column name format: %s", c.Name)
	}

	// 验证长度
	if c.Length <= 0 || c.Length > 255 {
		return fmt.Errorf("invalid varchar length: %d", c.Length)
	}

	return nil
}

// CreateDimReq 创建维度请求
type CreateDimReq struct {
	TableName     string         `json:"table_name"`
	DisplayName   string         `json:"display_name"`
	Description   string         `json:"description"`
	ParentID      uint           `json:"parent_id"`
	DimensionType string         `json:"dimension_type"`
	CustomColumns []CustomColumn `json:"custom_columns"` // 自定义列定义
}

// Validate 验证创建维度请求
func (r *CreateDimReq) Validate() error {
	if r.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if !utils.IsValidIdentifier(r.TableName) {
		return fmt.Errorf("invalid table name format: %s", r.TableName)
	}
	if r.DisplayName == "" {
		return fmt.Errorf("display name cannot be empty")
	}

	// 验证自定义列数量
	if len(r.CustomColumns) > 10 {
		return fmt.Errorf("too many custom columns, maximum is 10")
	}

	// 验证每个自定义列的定义
	columnNames := make(map[string]bool)
	for _, col := range r.CustomColumns {
		if err := col.Validate(); err != nil {
			return err
		}
		// 检查列名是否重复
		if columnNames[col.Name] {
			return fmt.Errorf("duplicate column name: %s", col.Name)
		}
		columnNames[col.Name] = true
	}

	return nil
}

// UpdateDimensionReq 更新维度请求
type UpdateDimensionReq struct {
	ID            uint           `json:"id"`
	TableName     string         `json:"table_name"`
	DisplayName   string         `json:"display_name"`
	Description   string         `json:"description"`
	CustomColumns []CustomColumn `json:"custom_columns"` // 自定义列定义
}

// Validate 验证更新维度请求
func (r *UpdateDimensionReq) Validate() error {
	if r.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if !utils.IsValidIdentifier(r.TableName) {
		return fmt.Errorf("invalid table name format: %s", r.TableName)
	}
	if r.DisplayName == "" {
		return fmt.Errorf("display name cannot be empty")
	}

	// 验证自定义列数量
	if len(r.CustomColumns) > 10 {
		return fmt.Errorf("too many custom columns, maximum is 10")
	}

	// 验证每个自定义列的定义
	columnNames := make(map[string]bool)
	for _, col := range r.CustomColumns {
		if err := col.Validate(); err != nil {
			return err
		}
		// 检查列名是否重复
		if columnNames[col.Name] {
			return fmt.Errorf("duplicate column name: %s", col.Name)
		}
		columnNames[col.Name] = true
	}

	return nil
}

type GetDimResp struct {
	ID            uint             `json:"id"`
	TableName     string           `json:"table_name"`
	DisplayName   string           `json:"display_name"`
	Description   string           `json:"description"`
	DimensionType string           `json:"dimension_type"`
	AppID         uint             `json:"app_id"`
	Status        int              `db:"status" json:"status"`
	CreatedAt     utils.CustomTime `db:"created_at" json:"created_at"`
	CreatorID     uint             `db:"creator_id" json:"creator_id"`
	UpdatedAt     utils.CustomTime `db:"updated_at" json:"updated_at"`
	UpdaterID     uint             `db:"updater_id" json:"updater_id"`
	CustomColumns []CustomColumn   `json:"custom_columns"` // 自定义列定义
}
