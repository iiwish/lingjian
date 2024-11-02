package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

type ApplicationService struct{}

// CreateApplication 创建应用
func (s *ApplicationService) CreateApplication(name, code, description string) error {
	_, err := model.DB.Exec(`
		INSERT INTO applications (name, code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, 1, ?, ?)
	`, name, code, description, time.Now(), time.Now())
	return err
}

// AssignApplicationToUser 为用户分配应用
func (s *ApplicationService) AssignApplicationToUser(userID, applicationID uint, isDefault bool) error {
	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 如果设置为默认应用，先将其他应用设置为非默认
	if isDefault {
		_, err = tx.Exec(`
			UPDATE user_applications 
			SET is_default = FALSE 
			WHERE user_id = ?
		`, userID)
		if err != nil {
			return err
		}
	}

	// 分配应用
	_, err = tx.Exec(`
		INSERT INTO user_applications (user_id, application_id, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE is_default = ?
	`, userID, applicationID, isDefault, time.Now(), time.Now(), isDefault)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserApplications 获取用户的所有应用
func (s *ApplicationService) GetUserApplications(userID uint) ([]model.Application, error) {
	var apps []model.Application
	err := model.DB.Select(&apps, `
		SELECT a.* FROM applications a
		INNER JOIN user_applications ua ON a.id = ua.application_id
		WHERE ua.user_id = ? AND a.status = 1
		ORDER BY ua.is_default DESC
	`, userID)
	return apps, err
}

// GetDefaultApplication 获取用户的默认应用
func (s *ApplicationService) GetDefaultApplication(userID uint) (*model.Application, error) {
	var app model.Application
	err := model.DB.Get(&app, `
		SELECT a.* FROM applications a
		INNER JOIN user_applications ua ON a.id = ua.application_id
		WHERE ua.user_id = ? AND ua.is_default = TRUE AND a.status = 1
	`, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, err
}

// CreateApplicationTemplate 创建应用模板
func (s *ApplicationService) CreateApplicationTemplate(name, description, configuration string, price float64, creatorID uint) error {
	_, err := model.DB.Exec(`
		INSERT INTO application_templates 
		(name, description, configuration, price, creator_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, name, description, configuration, price, creatorID, time.Now(), time.Now())
	return err
}

// ListApplicationTemplates 列出应用模板
func (s *ApplicationService) ListApplicationTemplates(status int) ([]model.ApplicationTemplate, error) {
	var templates []model.ApplicationTemplate
	err := model.DB.Select(&templates, `
		SELECT * FROM application_templates
		WHERE status = ?
		ORDER BY downloads DESC
	`, status)
	return templates, err
}

// CreateApplicationFromTemplate 从模板创建应用
func (s *ApplicationService) CreateApplicationFromTemplate(templateID uint, userID uint, appName, appCode string) error {
	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取模板信息
	var template model.ApplicationTemplate
	err = tx.Get(&template, "SELECT * FROM application_templates WHERE id = ?", templateID)
	if err != nil {
		return err
	}

	// 创建应用
	result, err := tx.Exec(`
		INSERT INTO applications (name, code, description, status, created_at, updated_at)
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
		INSERT INTO user_applications (user_id, application_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, userID, appID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	// 更新模板下载次数
	_, err = tx.Exec(`
		UPDATE application_templates 
		SET downloads = downloads + 1 
		WHERE id = ?
	`, templateID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// PublishTemplate 发布模板
func (s *ApplicationService) PublishTemplate(templateID uint) error {
	result, err := model.DB.Exec(`
		UPDATE application_templates 
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
