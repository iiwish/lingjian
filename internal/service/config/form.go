package config

import (
	"encoding/json"
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
func (s *FormService) CreateForm(form *model.ConfigForm, creatorID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 设置初始版本
	form.Version = 1
	form.Status = 1

	// 插入表单配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_forms (
			app_id, name, code, type, table_id,
			layout, fields, rules, events,
			status, version, created_at, updated_at
		) VALUES (
			:app_id, :name, :code, :type, :table_id,
			:layout, :fields, :rules, :events,
			:status, :version, NOW(), NOW()
		)
	`, form)
	if err != nil {
		return fmt.Errorf("insert sys_config_forms failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}
	form.ID = uint(id)

	// 创建版本记录
	version := &model.ConfigVersion{
		AppID:      form.AppID,
		ConfigType: "form",
		ConfigID:   form.ID,
		Version:    1,
		Content:    form.Fields, // 使用字段配置作为版本内容
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

// UpdateForm 更新表单配置
func (s *FormService) UpdateForm(form *model.ConfigForm, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前版本
	var currentVersion int
	err = tx.Get(&currentVersion, "SELECT version FROM sys_config_forms WHERE id = ?", form.ID)
	if err != nil {
		return fmt.Errorf("get current version failed: %v", err)
	}

	// 更新版本号
	form.Version = currentVersion + 1

	// 更新表单配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_forms SET 
			name = :name,
			code = :code,
			type = :type,
			table_id = :table_id,
			layout = :layout,
			fields = :fields,
			rules = :rules,
			events = :events,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, form)
	if err != nil {
		return fmt.Errorf("update sys_config_forms failed: %v", err)
	}

	// 创建新的版本记录
	version := &model.ConfigVersion{
		AppID:      form.AppID,
		ConfigType: "form",
		ConfigID:   form.ID,
		Version:    form.Version,
		Content:    form.Fields, // 使用字段配置作为版本内容
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

	// 软删除表单配置（将状态设置为0）
	_, err = tx.Exec("UPDATE sys_config_forms SET status = 0, updated_at = NOW() WHERE id = ?", id)
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
	err := s.db.Select(&forms, "SELECT * FROM sys_config_forms WHERE app_id = ? AND status = 1 ORDER BY id DESC", appID)
	if err != nil {
		return nil, fmt.Errorf("list forms from sys_config_forms failed: %v", err)
	}
	return forms, nil
}

// GetFormVersions 获取表单配置版本历史
func (s *FormService) GetFormVersions(id uint) ([]model.ConfigVersion, error) {
	var versions []model.ConfigVersion
	err := s.db.Select(&versions, `
		SELECT * FROM sys_config_versions 
		WHERE config_type = 'form' AND config_id = ? 
		ORDER BY version DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get form versions from sys_config_versions failed: %v", err)
	}
	return versions, nil
}

// RollbackForm 回滚表单配置到指定版本
func (s *FormService) RollbackForm(id uint, version int, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取指定版本的配置内容
	var targetVersion model.ConfigVersion
	err = tx.Get(&targetVersion, `
		SELECT * FROM sys_config_versions 
		WHERE config_type = 'form' AND config_id = ? AND version = ?
	`, id, version)
	if err != nil {
		return fmt.Errorf("get target version from sys_config_versions failed: %v", err)
	}

	// 获取当前表单配置
	var form model.ConfigForm
	err = tx.Get(&form, "SELECT * FROM sys_config_forms WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get current form from sys_config_forms failed: %v", err)
	}

	// 解析版本内容
	var fields []model.FormField
	if err := json.Unmarshal([]byte(targetVersion.Content), &fields); err != nil {
		return fmt.Errorf("unmarshal fields failed: %v", err)
	}

	// 更新表单配置
	form.Fields = targetVersion.Content
	form.Version = form.Version + 1

	// 更新表单配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_forms SET 
			fields = :fields,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, form)
	if err != nil {
		return fmt.Errorf("update sys_config_forms failed: %v", err)
	}

	// 创建新的版本记录
	newVersion := &model.ConfigVersion{
		AppID:      form.AppID,
		ConfigType: "form",
		ConfigID:   form.ID,
		Version:    form.Version,
		Content:    form.Fields,
		Comment:    fmt.Sprintf("Rollback to version %d", version),
		CreatorID:  updaterID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO sys_config_versions (
			app_id, config_type, config_id, version,
			content, comment, creator_id, created_at
		) VALUES (
			:app_id, :config_type, :config_id, :version,
			:content, :comment, :creator_id, NOW()
		)
	`, newVersion)
	if err != nil {
		return fmt.Errorf("insert sys_config_versions failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}
