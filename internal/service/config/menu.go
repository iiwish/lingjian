package config

import (
	"encoding/json"
	"fmt"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/jmoiron/sqlx"
)

// MenuService 菜单配置服务
type MenuService struct {
	db *sqlx.DB
}

// NewMenuService 创建菜单配置服务实例
func NewMenuService(db *sqlx.DB) *MenuService {
	return &MenuService{db: db}
}

// CreateMenu 创建菜单配置
func (s *MenuService) CreateMenu(menu *model.ConfigMenu, creatorID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 设置初始版本
	menu.Version = 1
	menu.Status = 1

	// 插入菜单配置
	result, err := tx.NamedExec(`
		INSERT INTO config_menus (
			app_id, parent_id, name, code, icon,
			path, component, sort, status, version,
			created_at, updated_at
		) VALUES (
			:app_id, :parent_id, :name, :code, :icon,
			:path, :component, :sort, :status, :version,
			NOW(), NOW()
		)
	`, menu)
	if err != nil {
		return fmt.Errorf("insert config_menus failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed: %v", err)
	}
	menu.ID = uint(id)

	// 创建版本记录
	menuContent, err := json.Marshal(menu)
	if err != nil {
		return fmt.Errorf("marshal menu failed: %v", err)
	}

	version := &model.ConfigVersion{
		AppID:      menu.AppID,
		ConfigType: "menu",
		ConfigID:   menu.ID,
		Version:    1,
		Content:    string(menuContent), // 使用完整菜单配置作为版本内容
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

// UpdateMenu 更新菜单配置
func (s *MenuService) UpdateMenu(menu *model.ConfigMenu, updaterID uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前版本
	var currentVersion int
	err = tx.Get(&currentVersion, "SELECT version FROM config_menus WHERE id = ?", menu.ID)
	if err != nil {
		return fmt.Errorf("get current version failed: %v", err)
	}

	// 更新版本号
	menu.Version = currentVersion + 1

	// 更新菜单配置
	_, err = tx.NamedExec(`
		UPDATE config_menus SET 
			parent_id = :parent_id,
			name = :name,
			code = :code,
			icon = :icon,
			path = :path,
			component = :component,
			sort = :sort,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, menu)
	if err != nil {
		return fmt.Errorf("update config_menus failed: %v", err)
	}

	// 创建新的版本记录
	menuContent, err := json.Marshal(menu)
	if err != nil {
		return fmt.Errorf("marshal menu failed: %v", err)
	}

	version := &model.ConfigVersion{
		AppID:      menu.AppID,
		ConfigType: "menu",
		ConfigID:   menu.ID,
		Version:    menu.Version,
		Content:    string(menuContent),
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

// GetMenu 获取菜单配置
func (s *MenuService) GetMenu(id uint) (*model.ConfigMenu, error) {
	var menu model.ConfigMenu
	err := s.db.Get(&menu, "SELECT * FROM config_menus WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("get menu failed: %v", err)
	}
	return &menu, nil
}

// DeleteMenu 删除菜单配置
func (s *MenuService) DeleteMenu(id uint) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 软删除菜单配置（将状态设置为0）
	_, err = tx.Exec("UPDATE config_menus SET status = 0, updated_at = NOW() WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete menu failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// ListMenus 获取菜单配置列表
func (s *MenuService) ListMenus(appID uint) ([]model.ConfigMenu, error) {
	var menus []model.ConfigMenu
	err := s.db.Select(&menus, `
		SELECT * FROM config_menus 
		WHERE app_id = ? AND status = 1 
		ORDER BY sort ASC, id ASC
	`, appID)
	if err != nil {
		return nil, fmt.Errorf("list menus failed: %v", err)
	}
	return menus, nil
}

// GetMenuVersions 获取菜单配置版本历史
func (s *MenuService) GetMenuVersions(id uint) ([]model.ConfigVersion, error) {
	var versions []model.ConfigVersion
	err := s.db.Select(&versions, `
		SELECT * FROM config_versions 
		WHERE config_type = 'menu' AND config_id = ? 
		ORDER BY version DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get menu versions failed: %v", err)
	}
	return versions, nil
}

// RollbackMenu 回滚菜单配置到指定版本
func (s *MenuService) RollbackMenu(id uint, version int, updaterID uint) error {
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
		WHERE config_type = 'menu' AND config_id = ? AND version = ?
	`, id, version)
	if err != nil {
		return fmt.Errorf("get target version failed: %v", err)
	}

	// 获取当前菜单配置
	var menu model.ConfigMenu
	err = tx.Get(&menu, "SELECT * FROM config_menus WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("get current menu failed: %v", err)
	}

	// 解析版本内容
	var targetMenu model.ConfigMenu
	if err := json.Unmarshal([]byte(targetVersion.Content), &targetMenu); err != nil {
		return fmt.Errorf("unmarshal menu failed: %v", err)
	}

	// 更新菜单配置（保留ID和AppID）
	targetMenu.ID = menu.ID
	targetMenu.AppID = menu.AppID
	targetMenu.Version = menu.Version + 1

	// 更新菜单配置
	_, err = tx.NamedExec(`
		UPDATE config_menus SET 
			parent_id = :parent_id,
			name = :name,
			code = :code,
			icon = :icon,
			path = :path,
			component = :component,
			sort = :sort,
			status = :status,
			version = :version,
			updated_at = NOW()
		WHERE id = :id
	`, targetMenu)
	if err != nil {
		return fmt.Errorf("update menu failed: %v", err)
	}

	// 创建新的版本记录
	menuContent, err := json.Marshal(targetMenu)
	if err != nil {
		return fmt.Errorf("marshal menu failed: %v", err)
	}

	newVersion := &model.ConfigVersion{
		AppID:      targetMenu.AppID,
		ConfigType: "menu",
		ConfigID:   targetMenu.ID,
		Version:    targetMenu.Version,
		Content:    string(menuContent),
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
