package service

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// AppService 应用服务
type AppService struct{}

// ListApps 获取所有应用列表
func (s *AppService) ListApps() (map[string]interface{}, error) {
	var apps []map[string]interface{}
	query := `
		SELECT * FROM apps 
		WHERE status = 1
		ORDER BY created_at DESC
	`
	err := model.DB.Select(&apps, query)
	if err != nil {
		return nil, err
	}

	// 返回带有items字段的响应
	return map[string]interface{}{
		"items": apps,
		"total": len(apps),
	}, nil
}

// CreateApp 创建应用
func (s *AppService) CreateApp(name, code, description string) (map[string]interface{}, error) {
	// 检查应用代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM apps WHERE code = ?", code)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("应用代码已存在")
	}

	now := time.Now()
	result, err := model.DB.Exec(`
		INSERT INTO apps (name, code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, name, code, description, 1, now, now)
	if err != nil {
		return nil, err
	}

	// 获取新创建的应用ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 返回新创建的应用信息
	return map[string]interface{}{
		"id":          uint(id),
		"name":        name,
		"code":        code,
		"description": description,
		"status":      1,
		"created_at":  now,
		"updated_at":  now,
	}, nil
}

// AssignAppToUser 为用户分配应用
func (s *AppService) AssignAppToUser(userID, appID uint, isDefault bool) error {
	// 检查用户和应用是否存在
	var userCount, appCount int
	err := model.DB.Get(&userCount, "SELECT COUNT(*) FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	if userCount == 0 {
		return errors.New("用户不存在")
	}

	err = model.DB.Get(&appCount, "SELECT COUNT(*) FROM apps WHERE id = ?", appID)
	if err != nil {
		return err
	}
	if appCount == 0 {
		return errors.New("应用不存在")
	}

	// 如果设置为默认应用，先取消其他默认应用
	if isDefault {
		_, err = model.DB.Exec(`
			UPDATE user_apps SET is_default = 0
			WHERE user_id = ? AND is_default = 1
		`, userID)
		if err != nil {
			return err
		}
	}

	// 分配应用给用户
	_, err = model.DB.Exec(`
		INSERT INTO user_apps (user_id, app_id, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE is_default = ?, updated_at = ?
	`, userID, appID, isDefault, time.Now(), time.Now(), isDefault, time.Now())

	return err
}

// GetUserApps 获取用户的应用列表
func (s *AppService) GetUserApps(userID uint) ([]map[string]interface{}, error) {
	var apps []map[string]interface{}
	query := `
		SELECT a.*, ua.is_default
		FROM apps a
		JOIN user_apps ua ON a.id = ua.app_id
		WHERE ua.user_id = ? AND a.status = 1
		ORDER BY ua.is_default DESC, a.created_at DESC
	`
	err := model.DB.Select(&apps, query, userID)
	return apps, err
}

// GetDefaultApp 获取用户的默认应用
func (s *AppService) GetDefaultApp(userID uint) (map[string]interface{}, error) {
	var app map[string]interface{}
	query := `
		SELECT a.*
		FROM apps a
		JOIN user_apps ua ON a.id = ua.app_id
		WHERE ua.user_id = ? AND ua.is_default = 1 AND a.status = 1
		LIMIT 1
	`
	err := model.DB.Get(&app, query, userID)
	return app, err
}

// CreateAppTemplate 创建应用模板
func (s *AppService) CreateAppTemplate(name, description, configuration string, price float64, creatorID uint) error {
	_, err := model.DB.Exec(`
		INSERT INTO app_templates (name, description, configuration, price, creator_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, name, description, configuration, price, creatorID, 0, time.Now(), time.Now())
	return err
}

// ListAppTemplates 获取应用模板列表
func (s *AppService) ListAppTemplates(status int) ([]map[string]interface{}, error) {
	var templates []map[string]interface{}
	query := `
		SELECT t.*, u.username as creator_name
		FROM app_templates t
		JOIN users u ON t.creator_id = u.id
		WHERE t.status = ?
		ORDER BY t.downloads DESC, t.created_at DESC
	`
	err := model.DB.Select(&templates, query, status)
	return templates, err
}

// CreateAppFromTemplate 从模板创建应用
func (s *AppService) CreateAppFromTemplate(templateID, userID uint, name, code string) error {
	// 检查模板是否存在且已上架
	var template struct {
		Configuration string
		Status        int
	}
	err := model.DB.Get(&template, "SELECT configuration, status FROM app_templates WHERE id = ?", templateID)
	if err != nil {
		return err
	}
	if template.Status != 1 {
		return errors.New("模板未上架")
	}

	// 创建应用
	appInfo, err := s.CreateApp(name, code, "从模板创建")
	if err != nil {
		return err
	}

	// 分配应用给用户
	err = s.AssignAppToUser(userID, appInfo["id"].(uint), false)
	if err != nil {
		return err
	}

	// 更新模板下载次数
	_, err = model.DB.Exec("UPDATE app_templates SET downloads = downloads + 1 WHERE id = ?", templateID)
	return err
}

// PublishTemplate 发布模板
func (s *AppService) PublishTemplate(templateID uint) error {
	_, err := model.DB.Exec("UPDATE app_templates SET status = 1 WHERE id = ?", templateID)
	return err
}
