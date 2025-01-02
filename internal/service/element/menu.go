package element

import (
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type MenuService struct {
	db *sqlx.DB
}

func NewMenuService(db *sqlx.DB) *MenuService {
	return &MenuService{db: db}
}

// GetMenuItems 获取菜单明细
func (s *MenuService) GetMenuItems(userID uint, menuID uint, itemID uint, queryType string, queryLevel uint) ([]*model.TreeMenuItem, error) {
	dimService := NewDimensionService(s.db)
	dimItems, err := dimService.TreeDimensionItems(userID, menuID, itemID, queryType, queryLevel)
	if err != nil {
		return nil, err
	}

	// 转换函数
	var convert func(items []model.TreeDimensionItem) []*model.TreeMenuItem
	convert = func(items []model.TreeDimensionItem) []*model.TreeMenuItem {
		var menuItems []*model.TreeMenuItem
		for _, item := range items {
			menuItem := &model.TreeMenuItem{
				ID:        item.ID,
				NodeID:    item.NodeID,
				ParentID:  item.ParentID,
				MenuName:  item.Name,
				MenuCode:  item.Code,
				Level:     item.Level,
				Sort:      item.Sort,
				Status:    item.Status,
				CreatedAt: item.CreatedAt,
				CreatorID: item.CreatorID,
				UpdatedAt: item.UpdatedAt,
				UpdaterID: item.UpdaterID,
				// 从CustomData中获取自定义列
				SourceID: utils.ParseUint(item.CustomData["source_id"]),
				MenuType: utils.ParseInt(item.CustomData["menu_type"]),
				IconPath: item.CustomData["icon_path"],
			}

			// 递归转换子节点
			if len(item.Children) > 0 {
				// 将[]*TreeDimensionItem转换为[]TreeDimensionItem
				children := make([]model.TreeDimensionItem, len(item.Children))
				for i, child := range item.Children {
					children[i] = *child
				}
				menuItem.Children = convert(children)
			}

			menuItems = append(menuItems, menuItem)
		}
		return menuItems
	}

	return convert(dimItems), nil
}

// CreateMenu 创建菜单明细
func (s *MenuService) CreateMenuItem(userID uint, menu *model.CreateMenuItemReq, menuID uint) (uint, error) {
	dimService := NewDimensionService(s.db)

	req := &model.CreateDimensionItemReq{
		Name:        menu.Name,
		Code:        menu.Code,
		Description: menu.Description,
		Status:      menu.Status,
		CustomData: map[string]string{
			"source_id": utils.Uint2String(menu.SourceID),
			"menu_type": utils.Int2String(menu.MenuType),
			"icon_path": menu.IconPath,
		},
		ParentID: menu.ParentID,
	}

	return dimService.CreateDimensionItem(req, userID, menuID)
}

// UpdateMenu 更新菜单明细
func (s *MenuService) UpdateMenuItem(menu *model.UpdateMenuItemReq, userID uint, menuID uint) error {
	dimService := NewDimensionService(s.db)

	req := &model.UpdateDimensionItemReq{
		ID:          menu.ID,
		Name:        menu.Name,
		Code:        menu.Code,
		Description: menu.Description,
		Status:      menu.Status,
		CustomData: map[string]string{
			"source_id": utils.Uint2String(menu.SourceID),
			"menu_type": utils.Int2String(menu.MenuType),
			"icon_path": menu.IconPath,
		},
	}

	return dimService.UpdateDimensionItem(req, userID, menuID)
}

// DeleteMenu 删除菜单明细
func (s *MenuService) DeleteMenuItem(operatorID uint, menuID uint, itemIDs []uint) error {
	dimService := NewDimensionService(s.db)
	return dimService.DeleteDimensionItems(operatorID, menuID, itemIDs)
}

// CreateSysMenu 创建系统菜单
func (s *MenuService) CreateSysMenu(appID uint, user_id uint, sysMenu *model.CreateMenuItemReq) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取系统菜单表名
	var sysMenuTableName string
	if appID == 0 {
		sysMenuTableName = "sys_menu_system"
	} else {
		sysMenuTableName = "app" + utils.Uint2String(appID) + "_menu_system"
	}

	// 获取系统菜单ID
	var sysMenuID uint
	err = tx.Get(&sysMenuID, "SELECT id FROM sys_config_dimensions WHERE table_name = ? AND app_id = ?", sysMenuTableName, appID)
	if err != nil {
		return err
	}

	// 插入系统菜单
	_, err = s.CreateMenuItem(user_id, sysMenu, sysMenuID)
	if err != nil {
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateMenuItemSort
func (s *MenuService) UpdateMenuItemSort(userID uint, menuID uint, itemID uint, parentID uint, sort int) error {
	dimService := NewDimensionService(s.db)
	return dimService.UpdateDimensionItemSort(userID, menuID, itemID, parentID, sort)
}
