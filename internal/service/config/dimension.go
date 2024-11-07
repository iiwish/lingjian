package config

import (
	"encoding/json"
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
func (s *DimensionService) CreateDimension(dimension *model.ConfigDimension, creatorID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入维度配置
	result, err := tx.NamedExec(`
		INSERT INTO config_dimensions (
			app_id, name, code, type, mysql_table_name,
			configuration, status, version, created_at, updated_at
		) VALUES (
			:app_id, :name, :code, :type, :mysql_table_name,
			:configuration, :status, :version, NOW(), NOW()
		)
	`, dimension)
	if err != nil {
		return fmt.Errorf("insert config_dimensions failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}

	// 创建版本记录
	version := &model.ConfigVersion{
		AppID:      dimension.AppID,
		ConfigType: "dimension",
		ConfigID:   uint(id),
		Version:    1,
		Content:    dimension.Configuration, // 使用配置作为版本内容
		CreatorID:  creatorID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO config_versions (
			app_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// UpdateDimension 更新维度配置
func (s *DimensionService) UpdateDimension(dimension *model.ConfigDimension, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前版本
	var currentVersion int
	err = tx.Get(&currentVersion, "SELECT version FROM config_dimensions WHERE id = ?", dimension.ID)
	if err != nil {
		return fmt.Errorf("get current version failed: %v", err)
	}

	// 更新版本号
	dimension.Version = currentVersion + 1

	// 更新维度配置
	_, err = tx.NamedExec(`
		UPDATE config_dimensions SET 
			name = :name,
			code = :code,
			type = :type,
			mysql_table_name = :mysql_table_name,
			configuration = :configuration,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, dimension)
	if err != nil {
		return fmt.Errorf("update config_dimensions failed: %v", err)
	}

	// 创建新的版本记录
	version := &model.ConfigVersion{
		AppID:      dimension.AppID,
		ConfigType: "dimension",
		ConfigID:   dimension.ID,
		Version:    dimension.Version,
		Content:    dimension.Configuration, // 使用配置作为版本内容
		CreatorID:  updaterID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO config_versions (
			app_id, config_type, config_id, version,
			content, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :creator_id, NOW()
		)
	`, version)
	if err != nil {
		return fmt.Errorf("insert config_versions failed: %v", err)
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
	err := s.db.Get(&dimension, "SELECT * FROM config_dimensions WHERE id = ?", id)
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
	_, err = tx.Exec("UPDATE config_dimensions SET status = 0, updated_at = NOW() WHERE id = ?", id)
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
	err := s.db.Select(&dimensions, "SELECT * FROM config_dimensions WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list dimensions failed: %v", err)
	}
	return dimensions, nil
}

// GetDimensionVersions 获取维度配置版本历史
func (s *DimensionService) GetDimensionVersions(id uint) ([]model.ConfigVersion, error) {
	var versions []model.ConfigVersion
	err := s.db.Select(&versions, `
		SELECT * FROM config_versions 
		WHERE config_type = 'dimension' AND config_id = ? 
		ORDER BY version DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get dimension versions failed: %v", err)
	}
	return versions, nil
}

// RollbackDimension 回滚维度配置到指定版本
func (s *DimensionService) RollbackDimension(id uint, version int, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取指定版本的配置内容
	var targetVersion model.ConfigVersion
	err = tx.Get(&targetVersion, `
		SELECT * FROM config_versions 
		WHERE config_type = 'dimension' AND config_id = ? AND version = ?
	`, id, version)
	if err != nil {
		return fmt.Errorf("get target version failed: %v", err)
	}

	// 获取当前维度配置
	var dimension model.ConfigDimension
	err = tx.Get(&dimension, "SELECT * FROM config_dimensions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get current dimension failed: %v", err)
	}

	// 解析版本内容
	var config model.DimensionConfig
	if err := json.Unmarshal([]byte(targetVersion.Content), &config); err != nil {
		return fmt.Errorf("unmarshal configuration failed: %v", err)
	}

	// 更新维度配置
	dimension.Configuration = targetVersion.Content
	dimension.Version = dimension.Version + 1

	// 更新维度配置
	_, err = tx.NamedExec(`
		UPDATE config_dimensions SET 
			configuration = :configuration,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, dimension)
	if err != nil {
		return fmt.Errorf("update dimension failed: %v", err)
	}

	// 创建新的版本记录
	newVersion := &model.ConfigVersion{
		AppID:      dimension.AppID,
		ConfigType: "dimension",
		ConfigID:   dimension.ID,
		Version:    dimension.Version,
		Content:    dimension.Configuration,
		Comment:    fmt.Sprintf("Rollback to version %d", version),
		CreatorID:  updaterID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO config_versions (
			app_id, config_type, config_id, version,
			content, comment, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :comment, :creator_id, NOW()
		)
	`, newVersion)
	if err != nil {
		return fmt.Errorf("insert version failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetDimensionValues 获取维度值列表
func (s *DimensionService) GetDimensionValues(dimensionID uint, filter map[string]any) ([]map[string]any, error) {
	// 获取维度配置
	dimension, err := s.GetDimension(dimensionID)
	if err != nil {
		return nil, err
	}

	// 解析维度配置
	var config model.DimensionConfig
	if err := json.Unmarshal([]byte(dimension.Configuration), &config); err != nil {
		return nil, fmt.Errorf("unmarshal configuration failed: %v", err)
	}

	// 构建查询SQL
	query := fmt.Sprintf("SELECT * FROM %s", dimension.MySQLTableName)
	var params []any

	// 添加过滤条件
	if len(filter) > 0 {
		query += " WHERE "
		var conditions []string
		for k, v := range filter {
			conditions = append(conditions, fmt.Sprintf("%s = ?", k))
			params = append(params, v)
		}
		query += fmt.Sprint(conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += fmt.Sprintf(" AND %s", conditions[i])
		}
	}

	// 添加排序
	if config.OrderField != "" {
		query += fmt.Sprintf(" ORDER BY %s", config.OrderField)
	}

	// 执行查询
	var values []map[string]any
	err = s.db.Select(&values, query, params...)
	if err != nil {
		return nil, fmt.Errorf("query dimension values failed: %v", err)
	}

	return values, nil
}
