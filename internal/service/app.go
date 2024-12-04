package service

import (
	"errors"
	"strconv"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/jmoiron/sqlx"
)

// AppService 应用服务
type AppService struct{}

// ListApps 获取所有应用列表
func (s *AppService) ListApps(userID uint) (map[string]interface{}, error) {
	var apps []model.App

	// 基础查询,添加权限控制
	appIDs := []uint{}
	query := `
        SELECT DISTINCT m.app_id FROM sys_config_menus m
        INNER JOIN sys_permissions p ON m.id = p.menu_id
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

	// 插入 'sys' 菜单
	sysMenu := &model.ConfigMenu{
		AppID:     uint(appID),
		ParentID:  0,
		NodeID:    "", // 先设置为空，避免唯一索引冲突
		MenuName:  "系统",
		MenuCode:  "_sys",
		MenuType:  1,
		SourceID:  0,
		Icon:      "folder",
		Sort:      1,
		Status:    1,
		CreatorID: user_id,
		UpdaterID: user_id,
	}

	// 插入 'sys' 菜单
	menuResult, err := tx.NamedExec(`
        INSERT INTO sys_config_menus (app_id, parent_id, node_id, menu_name, menu_code, menu_type, source_id, icon, sort, status, created_at, creator_id, updated_at, updater_id)
        VALUES (:app_id, :parent_id, :node_id, :menu_name, :menu_code, :menu_type, :source_id, :icon, :sort, :status, NOW(), :creator_id, NOW(), :updater_id)
    `, sysMenu)
	if err != nil {
		return nil, err
	}

	// 获取 'sys' 菜单的 ID
	sysMenuID, err := menuResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	sysMenu.ID = uint(sysMenuID)

	// 更新 'sys' 菜单的 NodeID
	sysMenu.NodeID = strconv.FormatUint(uint64(sysMenuID), 10)
	_, err = tx.Exec(`
        UPDATE sys_config_menus SET node_id = ? WHERE id = ?
    `, sysMenu.NodeID, sysMenuID)
	if err != nil {
		return nil, err
	}

	// 插入子菜单
	childMenus := []struct {
		MenuName string
		MenuCode string
		Icon     string
		Sort     int
	}{
		{"维度", "dimension", "folder", 1},
		{"表单", "form", "folder", 2},
		{"菜单", "menu", "folder", 3},
		{"模型", "model", "folder", 4},
		{"数据表", "table", "folder", 5},
	}

	for _, item := range childMenus {
		menu := &model.ConfigMenu{
			AppID:     uint(appID),
			ParentID:  sysMenu.ID,
			NodeID:    "", // 先设置为空，避免唯一索引冲突
			MenuName:  item.MenuName,
			MenuCode:  item.MenuCode,
			MenuType:  1,
			SourceID:  0,
			Icon:      item.Icon,
			Sort:      item.Sort,
			Status:    1,
			CreatorID: user_id,
			UpdaterID: user_id,
		}

		// 插入子菜单
		menuResult, err := tx.NamedExec(`
            INSERT INTO sys_config_menus (app_id, parent_id, node_id, menu_name, menu_code, menu_type, source_id, icon, sort, status, created_at, creator_id, updated_at, updater_id)
            VALUES (:app_id, :parent_id, :node_id, :menu_name, :menu_code, :menu_type, :source_id, :icon, :sort, :status, NOW(), :creator_id, NOW(), :updater_id)
        `, menu)
		if err != nil {
			return nil, err
		}

		// 获取子菜单的 ID
		menuID, err := menuResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		menu.ID = uint(menuID)

		// 更新子菜单的 NodeID
		menu.NodeID = sysMenu.NodeID + "_" + strconv.FormatUint(uint64(menuID), 10)
		_, err = tx.Exec(`
            UPDATE sys_config_menus SET node_id = ? WHERE id = ?
        `, menu.NodeID, menuID)
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
		MenuID:      sysMenu.ID,
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
func (s *AppService) UpdateApp(app *model.App, user_id uint) (*model.App, error) {
	app.UpdaterID = user_id

	// 检查应用代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_apps WHERE code = ? AND id != ?", app.Code, app.ID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("应用代码已存在")
	}

	// 更新应用信息
	_, err = model.DB.NamedExec(`
		UPDATE sys_apps SET name = :name, code = :code, description = :description, status = :status, updated_at = NOW(), updater_id = :updater_id
		WHERE id = :id
	`, app)
	if err != nil {
		return nil, err
	}

	return app, nil
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
	_, err = tx.Exec("DELETE FROM sys_config_menus WHERE app_id = ?", id)
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
