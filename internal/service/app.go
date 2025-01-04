package service

import (
	"errors"
	"strconv"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service/config"
	"github.com/iiwish/lingjian/internal/service/element"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/jmoiron/sqlx"
)

// AppService 应用服务
type AppService struct {
	menuService        *config.MenuService
	elementMenuService *element.MenuService
}

// NewAppService 创建应用服务实例
func NewAppService(db *sqlx.DB) *AppService {
	return &AppService{
		menuService:        config.NewMenuService(db),
		elementMenuService: element.NewMenuService(db),
	}
}

// ListApps 获取所有应用列表
func (s *AppService) ListApps(userID uint) (map[string]interface{}, error) {
	var apps []model.App

	// 基础查询,添加权限控制
	appIDs := []uint{}
	query := `
        SELECT DISTINCT m.app_id FROM sys_config_dimensions m
        INNER JOIN sys_permissions p ON m.id = p.dim_id
        INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
        INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ?
        AND m.status = 1
        AND p.status = 1
    `
	err := model.DB.Select(&appIDs, query, userID)
	if err != nil {
		return nil, err
	}

	if len(appIDs) == 0 {
		return map[string]interface{}{
			"items": []model.App{},
			"total": 0,
		}, nil
	}

	// 查询应用列表
	query = `
		SELECT * FROM sys_apps
		WHERE id IN (?)
	`
	query, args, err := sqlx.In(query, appIDs)
	if err != nil {
		return nil, err
	}
	err = model.DB.Select(&apps, query, args...)
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
func (s *AppService) CreateApp(app *model.App, user_id uint) (*model.App, error) {
	app.Status = 1
	app.CreatorID = user_id
	app.UpdaterID = user_id

	// 检查应用代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_apps WHERE code = ?", app.Code)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("应用代码已存在")
	}

	// 开启事务
	tx, err := model.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 插入应用
	result, err := tx.NamedExec(`
		INSERT INTO sys_apps (name, code, description, status, created_at, creator_id, updated_at, updater_id)
		VALUES (:name, :code, :description, :status, NOW(), :creator_id, NOW(), :updater_id)
	`, app)
	if err != nil {
		return nil, err
	}

	// 获取新创建的应用ID
	appID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 创建系统菜单
	sysMenuReq := &model.CreateMenuReq{
		TableName:   "app" + strconv.FormatInt(appID, 10) + "_menu_system",
		MenuName:    "系统",
		Description: "系统菜单",
		ParentID:    0,
	}

	sysMenuID, err := s.menuService.CreateMenu(user_id, uint(appID), sysMenuReq)
	if err != nil {
		return nil, err
	}

	// 创建子菜单项
	childMenus := []struct {
		Name     string
		Code     string
		IconPath string
	}{
		{"数据表", "table", "folder"},
		{"维度", "dimension", "folder"},
		{"菜单", "menu", "folder"},
		{"表单", "form", "folder"},
		{"模型", "model", "folder"},
	}

	for _, item := range childMenus {
		menuItemReq := &model.CreateMenuItemReq{
			MenuName: item.Name,
			MenuCode: item.Code,
			MenuType: 1,
			SourceID: 0,
			IconPath: item.IconPath,
			Status:   1,
			ParentID: 0,
		}

		_, err := s.elementMenuService.CreateMenuItem(user_id, menuItemReq, sysMenuID)
		if err != nil {
			return nil, err
		}
	}

	// 添加 'sys' 菜单的权限
	permission := &model.Permission{
		Name:        "系统菜单",
		Code:        "app" + strconv.FormatInt(appID, 10) + "_menu",
		Type:        "menu",
		Path:        "",
		Method:      "",
		DimID:       sysMenuID,
		ItemID:      0,
		Status:      1,
		Description: "app" + strconv.FormatInt(appID, 10) + "的系统菜单权限",
		CreatorID:   user_id,
		UpdaterID:   user_id,
	}

	// 插入权限
	permResult, err := tx.NamedExec(`
		INSERT INTO sys_permissions (name, code, type, path, method, menu_id, status, description, created_at, creator_id, updated_at, updater_id)
		VALUES (:name, :code, :type, :path, :method, :menu_id, :status, :description, NOW(), :creator_id, NOW(), :updater_id)
	`, permission)
	if err != nil {
		return nil, err
	}

	// 获取权限ID
	permID, err := permResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 为权限分配给角色
	rolePermission := &model.RolePermission{
		RoleID:       1, // 默认角色ID
		PermissionID: uint(permID),
		CreatorID:    user_id,
	}

	// 插入角色权限关联
	_, err = tx.NamedExec(`
		INSERT INTO sys_role_permissions (role_id, permission_id, created_at, creator_id)
		VALUES (:role_id, :permission_id, NOW(), :creator_id)
	`, rolePermission)
	if err != nil {
		return nil, err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 返回新创建的应用信息
	return &model.App{
		ID:          uint(appID),
		Name:        app.Name,
		Code:        app.Code,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   utils.NowCustomTime(),
		CreatorID:   user_id,
		UpdatedAt:   utils.NowCustomTime(),
		UpdaterID:   user_id,
	}, nil
}

// GetAppByID 根据 ID 获取应用信息
func (s *AppService) GetAppByID(id uint, user_id uint) (*model.App, error) {
	var app model.App
	err := model.DB.Get(&app, "SELECT * FROM sys_apps WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// UpdateApp 更新应用信息
func (s *AppService) UpdateApp(app *model.App, user_id uint) error {
	app.UpdaterID = user_id

	// 检查应用代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_apps WHERE code = ? AND id != ?", app.Code, app.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("应用代码已存在")
	}

	// 更新应用信息
	_, err = model.DB.NamedExec(`
		UPDATE sys_apps SET name = :name, code = :code, description = :description, status = :status, updated_at = NOW(), updater_id = :updater_id
		WHERE id = :id
	`, app)
	if err != nil {
		return err
	}

	return nil
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(id uint, user_id uint) error {
	// 开启事务
	tx, err := model.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除应用
	_, err = tx.Exec("DELETE FROM sys_apps WHERE id = ?", id)
	if err != nil {
		return err
	}

	// 删除应用下的菜单
	_, err = tx.Exec("DELETE FROM sys_config_dimensions WHERE app_id = ?", id)
	if err != nil {
		return err
	}

	// 删除应用下的角色权限关联
	_, err = tx.Exec("DELETE FROM sys_role_permissions WHERE permission_id IN (SELECT id FROM sys_permissions WHERE code LIKE ?)", "app"+strconv.FormatUint(uint64(id), 10)+"%")
	if err != nil {
		return err
	}

	// 删除应用下的权限
	_, err = tx.Exec("DELETE FROM sys_permissions WHERE code LIKE ?", "app"+strconv.FormatUint(uint64(id), 10)+"%")
	if err != nil {
		return err
	}

	// todo 考虑是否需要删除数据表

	// 提交事务
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
