package service

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
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
		INSERT INTO sys_apps (name, code, description, status, created_at, creator_id, updated_at, updated_id)
		VALUES (:name, :code, :description, :status, NOW(), :creator_id, NOW(), :updated_id)
	`, app)
	if err != nil {
		return nil, err
	}

	// 获取新创建的应用ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 插入默认菜单
	menus := []model.ConfigMenu{
		{AppID: uint(id), ParentID: 0, MenuName: "首页", MenuCode: "home", MenuType: "menu", Path: "/", Icon: "home", Sort: 0, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 0, MenuName: "文件夹", MenuCode: "folder", MenuType: "folder", Path: "", Icon: "folder", Sort: 1, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 2, MenuName: "系统", MenuCode: "sys", MenuType: "folder", Path: "", Icon: "settings", Sort: 1, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 3, MenuName: "维度", MenuCode: "dimension", MenuType: "menu", Path: "/sys/dimension", Icon: "dimension", Sort: 1, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 3, MenuName: "表单", MenuCode: "form", MenuType: "menu", Path: "/sys/form", Icon: "form", Sort: 2, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 3, MenuName: "菜单", MenuCode: "menu", MenuType: "menu", Path: "/sys/menu", Icon: "menu", Sort: 3, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 3, MenuName: "模型", MenuCode: "model", MenuType: "menu", Path: "/sys/model", Icon: "model", Sort: 4, Status: 1, CreatorID: user_id, UpdaterID: user_id},
		{AppID: uint(id), ParentID: 3, MenuName: "数据表", MenuCode: "table", MenuType: "menu", Path: "/sys/table", Icon: "table", Sort: 5, Status: 1, CreatorID: user_id, UpdaterID: user_id},
	}

	for _, menu := range menus {
		_, err = tx.NamedExec(`
            INSERT INTO sys_config_menus (app_id, parent_id, menu_name, menu_code, menu_type, path, icon, sort, status, created_at, creator_id, updated_at, updated_id)
            VALUES (:app_id, :parent_id, :menu_name, :menu_code, :menu_type, :path, :icon, :sort, :status, NOW(), :creator_id, NOW(), :updated_id)
        `, menu)
		if err != nil {
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 返回新创建的应用信息
	return &model.App{
		ID:          uint(id),
		Name:        app.Name,
		Code:        app.Code,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   time.Now(),
		CreatorID:   user_id,
		UpdatedAt:   time.Now(),
		UpdaterID:   user_id,
	}, nil
}
