package config

import (
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// FormService 表单配置服务
type FormService struct {
	db *sqlx.DB
}

// NewFormService 创建表单配置服务实例
func NewFormService(db *sqlx.DB) *FormService {
	return &FormService{db: db}
}

// CreateForm 创建表单配置
func (s *FormService) CreateForm(form *model.ConfigForm, creatorID uint) (uint, error) {
	form.Status = 1
	form.CreatorID = creatorID
	form.UpdaterID = creatorID

	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	form.Status = 1

	// 插入表单配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_forms (
			app_id, model_id, form_name, form_type, display_name, description, configuration, status, version, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :model_id, :form_name, :form_type, :display_name, :description, :configuration, :status, :version, NOW(), :creator_id, NOW(), :creator_id
		)
	`, form)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_forms failed: %v", err)
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

// UpdateForm 更新表单配置
func (s *FormService) UpdateForm(form *model.ConfigForm, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 更新表单配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_forms SET 
			app_id = :app_id,
			model_id = :model_id,
			form_name = :form_name,
			form_type = :form_type,
			display_name = :display_name,
			description = :description,
			configuration = :configuration,
			updated_at = NOW(),
			updater_id = :updater_id
		WHERE id = :id
	`, form)
	if err != nil {
		return fmt.Errorf("update sys_config_forms failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetForm 获取表单配置
func (s *FormService) GetForm(id uint) (*model.ConfigForm, error) {
	var form model.ConfigForm
	err := s.db.Get(&form, "SELECT * FROM sys_config_forms WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get form from sys_config_forms failed: %v", err)
	}
	return &form, nil
}

// DeleteForm 删除表单配置
func (s *FormService) DeleteForm(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 删除表单配置
	_, err = tx.Exec("DELETE FROM sys_config_forms WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete form from sys_config_forms failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListForms 获取表单配置列表
func (s *FormService) ListForms(appID uint) ([]model.ConfigForm, error) {
	var forms []model.ConfigForm
	err := s.db.Select(&forms, "SELECT id, app_id, model_id, form_name, form_type, display_name, description, status, version, created_at, creator_id, updated_at, updater_id FROM sys_config_forms WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list forms from sys_config_forms failed: %v", err)
	}
	return forms, nil
}
