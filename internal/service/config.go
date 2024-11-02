package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

type ConfigService struct{}

// CreateTableRequest 创建数据表配置请求
type CreateTableRequest struct {
	ApplicationID  uint         `json:"application_id" binding:"required"`
	Name           string       `json:"name" binding:"required"`
	Code           string       `json:"code" binding:"required"`
	Description    string       `json:"description"`
	MySQLTableName string       `json:"mysql_table_name" binding:"required"`
	Fields         []TableField `json:"fields" binding:"required"`
	Indexes        []TableIndex `json:"indexes"`
}

// TableField 表字段定义
type TableField struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"` // int, varchar, text, datetime等
	Length   int    `json:"length,omitempty"`
	Required bool   `json:"required"`
	Default  any    `json:"default,omitempty"`
}

// TableIndex 表索引定义
type TableIndex struct {
	Name   string   `json:"name" binding:"required"`
	Type   string   `json:"type" binding:"required"` // normal, unique
	Fields []string `json:"fields" binding:"required"`
}

// CreateTable 创建数据表配置
func (s *ConfigService) CreateTable(req *CreateTableRequest, creatorID uint) error {
	// 验证MySQL表是否存在
	var exists bool
	err := model.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = DATABASE() 
			AND table_name = ?
		)
	`, req.MySQLTableName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("MySQL表 %s 不存在", req.MySQLTableName)
	}

	// 验证字段是否与MySQL表匹配
	for _, field := range req.Fields {
		var columnExists bool
		err := model.DB.Get(&columnExists, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_schema = DATABASE() 
				AND table_name = ? 
				AND column_name = ?
			)
		`, req.MySQLTableName, field.Name)
		if err != nil {
			return err
		}
		if !columnExists {
			return fmt.Errorf("字段 %s 在MySQL表中不存在", field.Name)
		}
	}

	fieldsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		return err
	}

	indexesJSON, err := json.Marshal(req.Indexes)
	if err != nil {
		return err
	}

	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 创建表配置
	result, err := tx.Exec(`
		INSERT INTO config_tables (
			application_id, name, code, description, 
			mysql_table_name, fields, indexes, status, version, 
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 1, 1, ?, ?)
	`, req.ApplicationID, req.Name, req.Code, req.Description,
		req.MySQLTableName, fieldsJSON, indexesJSON,
		time.Now(), time.Now())
	if err != nil {
		return err
	}

	configID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 记录版本
	content := map[string]interface{}{
		"name":             req.Name,
		"code":             req.Code,
		"description":      req.Description,
		"mysql_table_name": req.MySQLTableName,
		"fields":           req.Fields,
		"indexes":          req.Indexes,
	}
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO config_versions (
			application_id, config_type, config_id,
			version, content, creator_id, created_at
		) VALUES (?, 'table', ?, 1, ?, ?, ?)
	`, req.ApplicationID, configID, contentJSON, creatorID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CreateDimensionRequest 创建维度配置请求
type CreateDimensionRequest struct {
	ApplicationID  uint            `json:"application_id" binding:"required"`
	Name           string          `json:"name" binding:"required"`
	Code           string          `json:"code" binding:"required"`
	Type           string          `json:"type" binding:"required"`
	MySQLTableName string          `json:"mysql_table_name" binding:"required"`
	Configuration  DimensionConfig `json:"configuration" binding:"required"`
}

// DimensionConfig 维度配置
type DimensionConfig struct {
	Fields       []string       `json:"fields" binding:"required"`        // MySQL表中的字段
	DisplayField string         `json:"display_field" binding:"required"` // 显示字段
	ValueField   string         `json:"value_field" binding:"required"`   // 值字段
	Filter       map[string]any `json:"filter,omitempty"`                 // 过滤条件
	Sort         []string       `json:"sort,omitempty"`                   // 排序字段
}

// CreateDimension 创建维度配置
func (s *ConfigService) CreateDimension(req *CreateDimensionRequest, creatorID uint) error {
	// 验证MySQL表是否存在
	var exists bool
	err := model.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = DATABASE() 
			AND table_name = ?
		)
	`, req.MySQLTableName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("MySQL表 %s 不存在", req.MySQLTableName)
	}

	// 验证字段是否与MySQL表匹配
	for _, field := range req.Configuration.Fields {
		var columnExists bool
		err := model.DB.Get(&columnExists, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_schema = DATABASE() 
				AND table_name = ? 
				AND column_name = ?
			)
		`, req.MySQLTableName, field)
		if err != nil {
			return err
		}
		if !columnExists {
			return fmt.Errorf("字段 %s 在MySQL表中不存在", field)
		}
	}

	configJSON, err := json.Marshal(req.Configuration)
	if err != nil {
		return err
	}

	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 创建维度配置
	result, err := tx.Exec(`
		INSERT INTO config_dimensions (
			application_id, name, code, type,
			mysql_table_name, configuration, status, version,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, 1, 1, ?, ?)
	`, req.ApplicationID, req.Name, req.Code, req.Type,
		req.MySQLTableName, configJSON,
		time.Now(), time.Now())
	if err != nil {
		return err
	}

	configID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 记录版本
	content := map[string]interface{}{
		"name":             req.Name,
		"code":             req.Code,
		"type":             req.Type,
		"mysql_table_name": req.MySQLTableName,
		"configuration":    req.Configuration,
	}
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO config_versions (
			application_id, config_type, config_id,
			version, content, creator_id, created_at
		) VALUES (?, 'dimension', ?, 1, ?, ?, ?)
	`, req.ApplicationID, configID, contentJSON, creatorID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetDimensionValues 获取维度值列表
func (s *ConfigService) GetDimensionValues(dimensionID uint, filter map[string]any) ([]map[string]any, error) {
	var dimension model.ConfigDimension
	err := model.DB.Get(&dimension, "SELECT * FROM config_dimensions WHERE id = ?", dimensionID)
	if err == sql.ErrNoRows {
		return nil, errors.New("维度不存在")
	}
	if err != nil {
		return nil, err
	}

	var config DimensionConfig
	err = json.Unmarshal([]byte(dimension.Configuration), &config)
	if err != nil {
		return nil, err
	}

	// 构建查询SQL
	query := fmt.Sprintf("SELECT %s FROM %s WHERE 1=1",
		strings.Join(config.Fields, ","),
		dimension.MySQLTableName,
	)

	// 添加过滤条件
	args := make([]interface{}, 0)
	if filter != nil {
		for field, value := range filter {
			query += " AND " + field + " = ?"
			args = append(args, value)
		}
	}

	// 添加配置中的过滤条件
	if config.Filter != nil {
		for field, value := range config.Filter {
			query += " AND " + field + " = ?"
			args = append(args, value)
		}
	}

	// 添加排序
	if len(config.Sort) > 0 {
		query += " ORDER BY " + strings.Join(config.Sort, ",")
	}

	// 执行查询
	var values []map[string]any
	err = model.DB.Select(&values, query, args...)
	return values, err
}
