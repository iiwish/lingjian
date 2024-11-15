package config

import (
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// DimensionService 维度配置服务
type DimensionService struct {
	db *sqlx.DB
}

// NewDimensionService 创建维度配置服务实例
func NewDimensionService(db *sqlx.DB) *DimensionService {
	return &DimensionService{db: db}
}

// CreateDimension 创建维度配置
func (s *DimensionService) CreateDimension(dimension *model.ConfigDimension, creatorID uint) (uint, error) {
	dimension.Status = 1
	dimension.CreatorID = creatorID
	dimension.UpdaterID = creatorID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 检查表名是否已存在
	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", dimension.TableName)
	if err != nil {
		return 0, fmt.Errorf("check table name failed: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("table name already exists")
	}

	// 插入维度配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_dimensions (
			app_id, table_name, display_name, description, status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :table_name, :display_name, :description, :status, NOW(), :creator_id, NOW(), :creator_id
		)
	`, dimension)

	if err != nil {
		return 0, fmt.Errorf("insert sys_config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建维度数据表
	tableName := dimension.TableName
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE %s (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
			node_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '节点ID',
			parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父节点ID',
			name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '名称',
			code VARCHAR(100) NOT NULL DEFAULT '' COMMENT '编码',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			creator_id INT NOT NULL DEFAULT 0 COMMENT '创建者ID',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			updater_id INT NOT NULL DEFAULT 0 COMMENT '更新者ID',
			UNIQUE KEY uk_code (code)
		)
	`, tableName)
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		return 0, fmt.Errorf("create table failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateDimension 更新维度配置
func (s *DimensionService) UpdateDimension(dimension *model.ConfigDimension, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取旧数据表名
	var oldTableName string
	err = tx.Get(&oldTableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", dimension.ID)
	if err != nil {
		return fmt.Errorf("get old table name failed: %v", err)
	}

	// 对比数据表名是否有变化
	if oldTableName != dimension.TableName {
		// 检查新表名是否已存在
		var count int
		err = tx.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", dimension.TableName)
		if err != nil {
			return fmt.Errorf("check table name failed: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("table name already exists")
		}

		// 修改数据表名
		_, err = tx.Exec("RENAME TABLE " + oldTableName + " TO " + dimension.TableName)
		if err != nil {
			return fmt.Errorf("rename table failed: %v", err)
		}
	}

	// 更新维度配置
	_, err = tx.NamedExec(`
			UPDATE sys_config_dimensions SET 
				table_name = :table_name,
				display_name = :display_name, 
				description = :description, 
				updated_at = NOW(), 
				updater_id = :updater_id
			WHERE id = :id
		`, dimension)
	if err != nil {
		return fmt.Errorf("update sys_config_dimensions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetDimension 获取维度配置
func (s *DimensionService) GetDimension(id uint) (*model.ConfigDimension, error) {
	var dimension model.ConfigDimension
	err := s.db.Get(&dimension, "SELECT * FROM sys_config_dimensions WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get dimension failed: %v", err)
	}
	return &dimension, nil
}

// DeleteDimension 删除维度配置
func (s *DimensionService) DeleteDimension(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 软删除维度配置（将状态设置为0）
	_, err = tx.Exec("UPDATE sys_config_dimensions SET status = 0, updated_at = NOW() WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete dimension failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListDimensions 获取维度配置列表
func (s *DimensionService) ListDimensions(appID uint) ([]model.ConfigDimension, error) {
	var dimensions []model.ConfigDimension
	err := s.db.Select(&dimensions, "SELECT * FROM sys_config_dimensions WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list dimensions failed: %v", err)
	}
	return dimensions, nil
}
