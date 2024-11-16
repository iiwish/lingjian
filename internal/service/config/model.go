package config

import (
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
func (s *ModelService) CreateModel(dataModel *model.ConfigModel, creatorID uint) (uint, error) {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 插入数据模型配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_models (
			app_id, model_name, display_name, discription, configuration,status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :model_name, :display_name, :discription, :configuration, :status, NOW(), :creator_id, NOW(), :creator_id
		)
	`, dataModel)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_models failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return uint(id), nil
}

// UpdateModel 更新数据模型配置
func (s *ModelService) UpdateModel(dataModel *model.ConfigModel, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 更新数据模型配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_models SET 
			app_id = :app_id,
			model_name = :model_name,
			display_name = :display_name,
			discription = :discription,
			configuration = :configuration,
			updated_at = NOW(),
			updater_id = :updater_id
		WHERE id = :id
	`, dataModel)
	if err != nil {
		return fmt.Errorf("update sys_config_models failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetModel 获取数据模型配置
func (s *ModelService) GetModel(id uint) (*model.ConfigModel, error) {
	var dataModel model.ConfigModel
	err := s.db.Get(&dataModel, "SELECT * FROM sys_config_models WHERE id = ?", id)
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

	// 删除数据模型配置
	_, err = tx.Exec("DELETE FROM sys_config_models WHERE id = ?", id)
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
func (s *ModelService) ListModels(appID uint) ([]model.ConfigModel, error) {
	var models []model.ConfigModel
	err := s.db.Select(&models, "SELECT id, app_id, model_name, display_name, description, status, version, created_at, creator_id, updated_at, updater_id FROM sys_config_models WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list models failed: %v", err)
	}
	return models, nil
}
