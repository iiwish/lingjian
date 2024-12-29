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
func (s *DimensionService) CreateDimension(req *model.CreateDimReq, creatorID uint, appID uint) (uint, error) {

	dimDB := model.ConfigDimension{
		AppID:       appID,
		TableName:   req.TableName,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      1,
		CreatorID:   creatorID,
		UpdaterID:   creatorID,
	}

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 检查表名是否已存在
	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", req.TableName)
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
	`, dimDB)

	if err != nil {
		return 0, fmt.Errorf("insert sys_config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建维度数据表
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE %s (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
			node_id VARCHAR(100) NOT NULL DEFAULT '' COMMENT '节点ID',
			parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父节点ID',
			name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '名称',
			code VARCHAR(100) NOT NULL DEFAULT '' COMMENT '编码',
			description VARCHAR(200) NOT NULL DEFAULT '' COMMENT '描述',
			level INT NOT NULL DEFAULT 0 COMMENT '层级',
			sort INT NOT NULL DEFAULT 0 COMMENT '排序',
			status TINYINT NOT NULL DEFAULT 1 COMMENT '状态',
			custom1 VARCHAR(100) NOT NULL DEFAULT '' COMMENT '自定义字段1',
			custom2 VARCHAR(100) NOT NULL DEFAULT '' COMMENT '自定义字段2',
			custom3 VARCHAR(100) NOT NULL DEFAULT '' COMMENT '自定义字段3',
			created_at DATETIME NOT NULL DEFAULT '1901-01-01 00:00:00' COMMENT '创建时间',
			creator_id INT NOT NULL DEFAULT 0 COMMENT '创建者ID',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			updater_id INT NOT NULL DEFAULT 0 COMMENT '更新者ID',
			UNIQUE KEY uk_code (code)
		)
	`, req.TableName)
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		return 0, fmt.Errorf("create table failed: %v", err)
	}

	// 插入菜单配置
	// 获取父节点node_id
	var parentNodeID string
	err = tx.Get(&parentNodeID, "SELECT node_id FROM sys_config_menus WHERE app_id = ? AND id = ?", appID, req.ParentID)
	if err != nil {
		return 0, fmt.Errorf("get parent node_id failed: %v", err)
	}
	// 获取sort
	var sort int
	err = tx.Get(&sort, "SELECT IFNULL(MAX(sort),0) FROM sys_config_menus WHERE app_id = ? AND parent_id = ?", appID, req.ParentID)
	if err != nil {
		return 0, fmt.Errorf("get sort failed: %v", err)
	}
	menuDB := model.ConfigMenu{
		AppID:     appID,
		NodeID:    parentNodeID + "_" + fmt.Sprint(id),
		ParentID:  req.ParentID,
		MenuName:  req.DisplayName,
		MenuCode:  fmt.Sprintf("dim_%d", id),
		MenuType:  3,
		Level:     1,
		Sort:      sort + 1,
		Icon:      "dimension",
		SourceID:  uint(id),
		CreatorID: creatorID,
		UpdaterID: creatorID,
	}
	_, err = tx.NamedExec(`
		INSERT INTO sys_config_menus (
			app_id, node_id, parent_id, menu_name, menu_code, menu_type, level, sort, icon, source_id, status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :node_id, :parent_id, :menu_name, :menu_code, :menu_type, :level, :sort, :icon, :source_id, 1, NOW(), :creator_id, NOW(), :creator_id
		)`, menuDB)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_menus failed: %v", err)
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

	// 获取数据表名
	var tableName string
	err = tx.Get(&tableName, "SELECT table_name FROM sys_config_dimensions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get table name failed: %v", err)
	}

	// 删除维度配置
	_, err = tx.Exec("DELETE FROM sys_config_dimensions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete dimension failed: %v", err)
	}

	// 删除数据表
	_, err = tx.Exec("DROP TABLE " + tableName)
	if err != nil {
		return fmt.Errorf("drop table failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
