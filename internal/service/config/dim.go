package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/element"
	"github.com/jmoiron/sqlx"
)

func toJSONString(v interface{}) string {
	if v == nil {
		return "[]"
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

// DimensionService 维度配置服务
type DimensionService struct {
	db          *sqlx.DB
	menuService *element.MenuService
}

// NewDimensionService 创建维度配置服务实例
func NewDimensionService(db *sqlx.DB) *DimensionService {
	return &DimensionService{
		db:          db,
		menuService: element.NewMenuService(db),
	}
}

// CreateDimension 创建维度配置
func (s *DimensionService) CreateDimension(req *model.CreateDimReq, creatorID uint, appID uint) (uint, error) {
	dimType := "general"
	if req.DimensionType != "" {
		dimType = req.DimensionType
	}
	// 维度配置
	dimDB := model.ConfigDimension{
		AppID:         appID,
		TableName:     req.TableName,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Status:        1,
		DimensionType: dimType,
		CustomColumns: toJSONString(req.CustomColumns),
		CreatorID:     creatorID,
		UpdaterID:     creatorID,
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
			app_id, table_name, display_name, description, dimension_type, status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :table_name, :display_name, :description, :dimension_type, :status, NOW(), :creator_id, NOW(), :creator_id
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
	var createTableSQL strings.Builder
	createTableSQL.WriteString(fmt.Sprintf(`
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
	`, req.TableName))

	// 添加自定义列
	for _, col := range req.CustomColumns {
		colDef := fmt.Sprintf("%s VARCHAR(%d) NOT NULL DEFAULT '' COMMENT '%s'",
			col.Name, col.Length, col.Comment)

		createTableSQL.WriteString(colDef + ",\n")
	}

	// 添加基础字段
	createTableSQL.WriteString(`
			created_at DATETIME NOT NULL DEFAULT '1901-01-01 00:00:00' COMMENT '创建时间',
			creator_id INT NOT NULL DEFAULT 0 COMMENT '创建者ID',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			updater_id INT NOT NULL DEFAULT 0 COMMENT '更新者ID',
			UNIQUE KEY uk_code (code)
		)
	`)
	_, err = tx.Exec(createTableSQL.String())
	if err != nil {
		return 0, fmt.Errorf("create table failed: %v", err)
	}

	if req.ParentID != 0 {
		var menuType int
		if req.DimensionType == "menu" {
			menuType = 4
		} else {
			menuType = 3
		}
		// 插入菜单配置
		menu := &model.CreateMenuItemReq{
			ParentID:    req.ParentID,
			Name:        req.DisplayName,
			Code:        req.TableName,
			Description: req.Description,
			Status:      1,
			SourceID:    uint(id),
			MenuType:    menuType,
			IconPath:    "folder",
		}
		_, err = s.menuService.CreateMenuItem(creatorID, menu, appID)
		if err != nil {
			return 0, fmt.Errorf("create menu item failed: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateDimension 更新维度配置
func (s *DimensionService) UpdateDimension(req *model.UpdateDimensionReq, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取旧数据表名和配置
	var oldDim model.ConfigDimension
	err = tx.Get(&oldDim, "SELECT * FROM sys_config_dimensions WHERE id = ?", req.ID)
	if err != nil {
		return fmt.Errorf("get old dimension failed: %v", err)
	}

	// 对比数据表名是否有变化
	if oldDim.TableName != req.TableName {
		// 检查新表名是否已存在
		var count int
		err = tx.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", req.TableName)
		if err != nil {
			return fmt.Errorf("check table name failed: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("table name already exists")
		}

		// 修改数据表名
		_, err = tx.Exec("RENAME TABLE " + oldDim.TableName + " TO " + req.TableName)
		if err != nil {
			return fmt.Errorf("rename table failed: %v", err)
		}
	}

	// 获取表的列信息
	var columns []struct {
		ColumnName string `db:"COLUMN_NAME"`
	}
	err = tx.Select(&columns, `
		SELECT COLUMN_NAME 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = ? 
		AND COLUMN_NAME NOT IN (
			'id', 'node_id', 'parent_id', 'name', 'code', 'description',
			'level', 'sort', 'status', 'created_at', 'creator_id', 
			'updated_at', 'updater_id'
		)
	`, req.TableName)
	if err != nil {
		return fmt.Errorf("get columns failed: %v", err)
	}

	// 构建当前自定义列map
	currentColumns := make(map[string]bool)
	for _, col := range columns {
		currentColumns[col.ColumnName] = true
	}

	// 构建新自定义列map
	newColumns := make(map[string]model.CustomColumn)
	for _, col := range req.CustomColumns {
		newColumns[col.Name] = col
	}

	// 删除不再需要的列
	for colName := range currentColumns {
		if _, exists := newColumns[colName]; !exists {
			_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", req.TableName, colName))
			if err != nil {
				return fmt.Errorf("drop column failed: %v", err)
			}
		}
	}

	// 添加新列
	for colName, col := range newColumns {
		if !currentColumns[colName] {
			colDef := fmt.Sprintf("ADD COLUMN %s VARCHAR(%d) NOT NULL DEFAULT '' COMMENT '%s'",
				col.Name, col.Length, col.Comment)

			_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s %s", req.TableName, colDef))
			if err != nil {
				return fmt.Errorf("add column failed: %v", err)
			}
		}
	}

	// 更新维度配置
	_, err = tx.Exec(`
		UPDATE sys_config_dimensions SET 
			table_name = ?,
			display_name = ?, 
			description = ?, 
			custom_columns = ?,
			updated_at = NOW(), 
			updater_id = ?
		WHERE id = ?
	`, req.TableName, req.DisplayName, req.Description, toJSONString(req.CustomColumns), updaterID, req.ID)
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
func (s *DimensionService) GetDimension(id uint) (*model.GetDimResp, error) {
	var dimension model.ConfigDimension
	err := s.db.Get(&dimension, "SELECT * FROM sys_config_dimensions WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get dimension failed: %v", err)
	}

	result := model.GetDimResp{
		ID:            dimension.ID,
		AppID:         dimension.AppID,
		TableName:     dimension.TableName,
		DisplayName:   dimension.DisplayName,
		Description:   dimension.Description,
		DimensionType: dimension.DimensionType,
		Status:        dimension.Status,
		CreatedAt:     dimension.CreatedAt,
		CreatorID:     dimension.CreatorID,
		UpdatedAt:     dimension.UpdatedAt,
		UpdaterID:     dimension.UpdaterID,
		CustomColumns: []model.CustomColumn{},
	}
	if err := json.Unmarshal([]byte(dimension.CustomColumns), &result.CustomColumns); err != nil {
		return nil, fmt.Errorf("unmarshal custom columns failed: %v", err)
	}

	return &result, nil
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

// GetDimensions 获取维度配置列表
func (s *DimensionService) GetDimensions(userID uint, appID uint, dimType string) ([]model.GetDimResp, error) {
	// 查询用户在该维度的权限
	permissionQuery := `
		SELECT DISTINCT p.dim_id FROM sys_permissions p
		INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
		INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? 
		AND p.status = 1
	`
	var permDimIDs []uint
	err := s.db.Select(&permDimIDs, permissionQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions failed: %v", err)
	}

	// 如果没有权限记录,直接返回空值
	if len(permDimIDs) == 0 {
		return []model.GetDimResp{}, nil
	}
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 查询维度配置
	var dimensions []model.ConfigDimension
	if dimType == "" {
		query, args, err := sqlx.In("SELECT * FROM sys_config_dimensions WHERE app_id = ? AND id IN (?)", appID, permDimIDs)
		if err != nil {
			return nil, fmt.Errorf("prepare query failed: %v", err)
		}
		query = tx.Rebind(query)
		err = tx.Select(&dimensions, query, args...)
		if err != nil {
			return nil, fmt.Errorf("get dimensions failed: %v", err)
		}
	} else {
		query, args, err := sqlx.In("SELECT * FROM sys_config_dimensions WHERE app_id = ? AND dimension_type = ? AND id IN (?)", appID, dimType, permDimIDs)
		if err != nil {
			return nil, fmt.Errorf("prepare query failed: %v", err)
		}
		query = tx.Rebind(query)
		err = tx.Select(&dimensions, query, args...)
		if err != nil {
			return nil, fmt.Errorf("get dimensions failed: %v", err)
		}
	}

	// 构建返回结果
	results := make([]model.GetDimResp, 0, len(dimensions))
	for _, dimension := range dimensions {
		result := model.GetDimResp{
			ID:            dimension.ID,
			AppID:         dimension.AppID,
			TableName:     dimension.TableName,
			DisplayName:   dimension.DisplayName,
			Description:   dimension.Description,
			DimensionType: dimension.DimensionType,
			Status:        dimension.Status,
			CreatedAt:     dimension.CreatedAt,
			CreatorID:     dimension.CreatorID,
			UpdatedAt:     dimension.UpdatedAt,
			UpdaterID:     dimension.UpdaterID,
			CustomColumns: []model.CustomColumn{},
		}
		if err := json.Unmarshal([]byte(dimension.CustomColumns), &result.CustomColumns); err != nil {
			return nil, fmt.Errorf("unmarshal custom columns failed: %v", err)
		}
		results = append(results, result)
	}

	return results, nil
}
