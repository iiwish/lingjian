package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

type AppService struct{}

// CreateApp 创建应用
func (s *AppService) CreateApp(name, code, description string) error {
	_, err := model.DB.Exec(`
		INSERT INTO apps (name, code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, 1, ?, ?)
	`, name, code, description, time.Now(), time.Now())
	return err
}

// AssignAppToUser 为用户分配应用
func (s *AppService) AssignAppToUser(userID, appID uint, isDefault bool) error {
	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 如果设置为默认应用，先将其他应用设置为非默认
	if isDefault {
		_, err = tx.Exec(`
			UPDATE user_apps 
			SET is_default = FALSE 
			WHERE user_id = ?
		`, userID)
		if err != nil {
			return err
		}
	}

	// 分配应用
	_, err = tx.Exec(`
		INSERT INTO user_apps (user_id, app_id, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE is_default = ?
	`, userID, appID, isDefault, time.Now(), time.Now(), isDefault)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserApps 获取用户的所有应用
func (s *AppService) GetUserApps(userID uint) ([]model.App, error) {
	var apps []model.App
	err := model.DB.Select(&apps, `
		SELECT a.* FROM apps a
		INNER JOIN user_apps ua ON a.id = ua.app_id
		WHERE ua.user_id = ? AND a.status = 1
		ORDER BY ua.is_default DESC
	`, userID)
	return apps, err
}

// GetDefaultApp 获取用户的默认应用
func (s *AppService) GetDefaultApp(userID uint) (*model.App, error) {
	var app model.App
	err := model.DB.Get(&app, `
		SELECT a.* FROM apps a
		INNER JOIN user_apps ua ON a.id = ua.app_id
		WHERE ua.user_id = ? AND ua.is_default = TRUE AND a.status = 1
	`, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, err
}

// CreateAppTemplate 创建应用模板
func (s *AppService) CreateAppTemplate(name, description, configuration string, price float64, creatorID uint) error {
	_, err := model.DB.Exec(`
		INSERT INTO app_templates 
		(name, description, configuration, price, creator_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, name, description, configuration, price, creatorID, time.Now(), time.Now())
	return err
}

// ListAppTemplates 列出应用模板
func (s *AppService) ListAppTemplates(status int) ([]model.AppTemplate, error) {
	var templates []model.AppTemplate
	err := model.DB.Select(&templates, `
		SELECT * FROM app_templates
		WHERE status = ?
		ORDER BY downloads DESC
	`, status)
	return templates, err
}

// CreateAppFromTemplate 从模板创建应用
func (s *AppService) CreateAppFromTemplate(templateID uint, userID uint, appName, appCode string) error {
	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取模板信息
	var template model.AppTemplate
	err = tx.Get(&template, "SELECT * FROM app_templates WHERE id = ?", templateID)
	if err != nil {
		return err
	}

	// 创建应用
	result, err := tx.Exec(`
		INSERT INTO apps (name, code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, 1, ?, ?)
	`, appName, appCode, template.Description, time.Now(), time.Now())
	if err != nil {
		return err
	}

	appID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 分配应用给用户
	_, err = tx.Exec(`
		INSERT INTO user_apps (user_id, app_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, userID, appID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	// 更新模板下载次数
	_, err = tx.Exec(`
		UPDATE app_templates 
		SET downloads = downloads + 1 
		WHERE id = ?
	`, templateID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// PublishTemplate 发布模板
func (s *AppService) PublishTemplate(templateID uint) error {
	result, err := model.DB.Exec(`
		UPDATE app_templates 
		SET status = 1 
		WHERE id = ?
	`, templateID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("模板不存在")
	}

	return nil
}
