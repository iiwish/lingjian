package config

import (
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// TableService 数据表配置服务
type TableService struct {
	db *sqlx.DB
}

// NewTableService 创建数据表配置服务实例
func NewTableService(db *sqlx.DB) *TableService {
	return &TableService{db: db}
}

// CreateTable 创建数据表配置
func (s *TableService) CreateTable(table *model.ConfigTable, creatorID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入数据表配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_tables (
			app_id, name, code, description, mysql_table_name,
			fields, indexes, status, version, created_at, updated_at
		) VALUES (
			:app_id, :name, :code, :description, :mysql_table_name,
			:fields, :indexes, :status, :version, NOW(), NOW()
		)
	`, table)
	if err != nil {
		return fmt.Errorf("insert sys_config_tables failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建版本记录
	version := &model.ConfigVersion{
		AppID:      table.AppID,
		ConfigType: "table",
		ConfigID:   uint(id),
		Version:    1,
		CreatorID:  creatorID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO sys_config_versions (
			app_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert sys_config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// UpdateTable 更新数据表配置
func (s *TableService) UpdateTable(table *model.ConfigTable, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前版本
	var currentVersion int
	err = tx.Get(&currentVersion, "SELECT version FROM sys_config_tables WHERE id = ?", table.ID)
	if err != nil {
		return fmt.Errorf("get current version failed: %v", err)
	}

	// 更新版本号
	table.Version = currentVersion + 1

	// 更新数据表配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_tables SET 
			name = :name,
			code = :code,
			description = :description,
			mysql_table_name = :mysql_table_name,
			fields = :fields,
			indexes = :indexes,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, table)
	if err != nil {
		return fmt.Errorf("update sys_config_tables failed: %v", err)
	}

	// 创建新的版本记录
	version := &model.ConfigVersion{
		AppID:      table.AppID,
		ConfigType: "table",
		ConfigID:   table.ID,
		Version:    table.Version,
		CreatorID:  updaterID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO sys_config_versions (
			app_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert sys_config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetTable 获取数据表配置
func (s *TableService) GetTable(id uint) (*model.ConfigTable, error) {
	var table model.ConfigTable
	err := s.db.Get(&table, "SELECT * FROM sys_config_tables WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get table failed: %v", err)
	}
	return &table, nil
}

// DeleteTable 删除数据表配置
func (s *TableService) DeleteTable(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 软删除数据表配置（将状态设置为0）
	_, err = tx.Exec("UPDATE sys_config_tables SET status = 0, updated_at = NOW() WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete table failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListTables 获取数据表配置列表
func (s *TableService) ListTables(appID uint) ([]model.ConfigTable, error) {
	var tables []model.ConfigTable
	err := s.db.Select(&tables, "SELECT * FROM sys_config_tables WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list tables failed: %v", err)
	}
	return tables, nil
}
