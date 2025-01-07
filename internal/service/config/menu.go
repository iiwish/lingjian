package config

import (
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
func (s *MenuService) CreateMenu(userID uint, appID uint, menu *model.CreateMenuReq) (uint, error) {
	dimService := NewDimensionService(s.db)

	req := &model.CreateDimReq{
		TableName:     menu.TableName,
		DisplayName:   menu.MenuName,
		Description:   menu.Description,
		DimensionType: "menu",
		ParentID:      menu.ParentID,
		CustomColumns: []model.CustomColumn{
			{Name: "source_id", Length: 30, Comment: "数据源id"},
			{Name: "menu_type", Length: 1, Comment: "菜单类型"},
			{Name: "icon_path", Length: 100, Comment: "图标路径"},
		},
	}

	return dimService.CreateDimension(req, userID, appID)
}

// UpdateMenu 更新菜单配置
func (s *MenuService) UpdateMenu(menu *model.UpdateMenuReq, userID uint, dimID uint) error {
	dimService := NewDimensionService(s.db)

	req := &model.UpdateDimensionReq{
		ID:          dimID,
		TableName:   menu.TableName,
		DisplayName: menu.MenuName,
		Description: "",
		CustomColumns: []model.CustomColumn{
			{Name: "source_id", Length: 30, Comment: "数据源id"},
			{Name: "menu_type", Length: 1, Comment: "菜单类型"},
			{Name: "icon_path", Length: 100, Comment: "图标路径"},
		},
	}

	return dimService.UpdateDimension(req, userID)
}

// DeleteMenu 删除菜单配置
func (s *MenuService) DeleteMenu(id uint) error {
	dimService := NewDimensionService(s.db)
	return dimService.DeleteDimension(id)
}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList(userID uint, appID uint) ([]model.GetDimResp, error) {
	dimService := NewDimensionService(s.db)
	return dimService.GetDimensions(userID, appID, "menu")
}

// GetMenuByID 获取菜单配置
func (s *MenuService) GetMenuByID(id uint) (*model.GetDimResp, error) {
	dimService := NewDimensionService(s.db)
	return dimService.GetDimension(id)
}

// GetSystemMenuID 获取系统菜单ID
func (s *MenuService) GetSystemMenuID(appID uint) (uint, error) {
	var id uint
	err := s.db.Get(&id, "SELECT id FROM sys_config_dimensions WHERE table_name like '%_menu_system' AND app_id = ?", appID)
	if err != nil {
		return 0, fmt.Errorf("get system menu id failed: %v", err)
	}
	return id, nil
}
