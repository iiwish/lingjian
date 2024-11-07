package config

import (
	"encoding/json"
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// ModelService 数据模型配置服务
type ModelService struct {
	db *sqlx.DB
}

// NewModelService 创建数据模型配置服务实例
func NewModelService(db *sqlx.DB) *ModelService {
	return &ModelService{db: db}
}

// CreateModel 创建数据模型配置
func (s *ModelService) CreateModel(dataModel *model.ConfigDataModel, creatorID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 设置初始版本
	dataModel.Version = 1
	dataModel.Status = 1

	// 插入数据模型配置
	result, err := tx.NamedExec(`
		INSERT INTO config_data_models (
			app_id, name, code, table_id, fields,
			dimensions, metrics, status, version,
			created_at, updated_at
		) VALUES (
			:app_id, :name, :code, :table_id, :fields,
			:dimensions, :metrics, :status, :version,
			NOW(), NOW()
		)
	`, dataModel)
	if err != nil {
		return fmt.Errorf("insert config_data_models failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}
	dataModel.ID = uint(id)

	// 创建版本记录
	version := &model.ConfigVersion{
		AppID:      dataModel.AppID,
		ConfigType: "model",
		ConfigID:   dataModel.ID,
		Version:    1,
		Content:    dataModel.Fields, // 使用字段配置作为版本内容
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

// UpdateModel 更新数据模型配置
func (s *ModelService) UpdateModel(dataModel *model.ConfigDataModel, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前版本
	var currentVersion int
	err = tx.Get(&currentVersion, "SELECT version FROM config_data_models WHERE id = ?", dataModel.ID)
	if err != nil {
		return fmt.Errorf("get current version failed: %v", err)
	}

	// 更新版本号
	dataModel.Version = currentVersion + 1

	// 更新数据模型配置
	_, err = tx.NamedExec(`
		UPDATE config_data_models SET 
			name = :name,
			code = :code,
			table_id = :table_id,
			fields = :fields,
			dimensions = :dimensions,
			metrics = :metrics,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, dataModel)
	if err != nil {
		return fmt.Errorf("update config_data_models failed: %v", err)
	}

	// 创建新的版本记录
	version := &model.ConfigVersion{
		AppID:      dataModel.AppID,
		ConfigType: "model",
		ConfigID:   dataModel.ID,
		Version:    dataModel.Version,
		Content:    dataModel.Fields, // 使用字段配置作为版本内容
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

// GetModel 获取数据模型配置
func (s *ModelService) GetModel(id uint) (*model.ConfigDataModel, error) {
	var dataModel model.ConfigDataModel
	err := s.db.Get(&dataModel, "SELECT * FROM config_data_models WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get model failed: %v", err)
	}
	return &dataModel, nil
}

// DeleteModel 删除数据模型配置
func (s *ModelService) DeleteModel(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 软删除数据模型配置（将状态设置为0）
	_, err = tx.Exec("UPDATE config_data_models SET status = 0, updated_at = NOW() WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete model failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListModels 获取数据模型配置列表
func (s *ModelService) ListModels(appID uint) ([]model.ConfigDataModel, error) {
	var models []model.ConfigDataModel
	err := s.db.Select(&models, "SELECT * FROM config_data_models WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list models failed: %v", err)
	}
	return models, nil
}

// GetModelVersions 获取数据模型配置版本历史
func (s *ModelService) GetModelVersions(id uint) ([]model.ConfigVersion, error) {
	var versions []model.ConfigVersion
	err := s.db.Select(&versions, `
		SELECT * FROM config_versions 
		WHERE config_type = 'model' AND config_id = ? 
		ORDER BY version DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get model versions failed: %v", err)
	}
	return versions, nil
}

// RollbackModel 回滚数据模型配置到指定版本
func (s *ModelService) RollbackModel(id uint, version int, updaterID uint) error {
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
		WHERE config_type = 'model' AND config_id = ? AND version = ?
	`, id, version)
	if err != nil {
		return fmt.Errorf("get target version failed: %v", err)
	}

	// 获取当前数据模型配置
	var dataModel model.ConfigDataModel
	err = tx.Get(&dataModel, "SELECT * FROM config_data_models WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get current model failed: %v", err)
	}

	// 解析版本内容
	var fields []model.ModelField
	if err := json.Unmarshal([]byte(targetVersion.Content), &fields); err != nil {
		return fmt.Errorf("unmarshal fields failed: %v", err)
	}

	// 更新数据模型配置
	dataModel.Fields = targetVersion.Content
	dataModel.Version = dataModel.Version + 1

	// 更新数据模型配置
	_, err = tx.NamedExec(`
		UPDATE config_data_models SET 
			fields = :fields,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, dataModel)
	if err != nil {
		return fmt.Errorf("update model failed: %v", err)
	}

	// 创建新的版本记录
	newVersion := &model.ConfigVersion{
		AppID:      dataModel.AppID,
		ConfigType: "model",
		ConfigID:   dataModel.ID,
		Version:    dataModel.Version,
		Content:    dataModel.Fields,
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
