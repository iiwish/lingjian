package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
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
func (s *MenuService) CreateMenu(menu *model.ConfigMenu, creatorID uint) (uint, error) {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	menu.Status = 1

	// 插入菜单配置
	result, err := tx.NamedExec(`
		INSERT INTO sys_config_menus (
			app_id, parent_id, node_id, menu_name, menu_code, menu_type, level, sort, icon, path, status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :parent_id, :node_id, :menu_name, :menu_code, :menu_type, :level, :sort, :icon, :path, :status, NOW(), :creator_id, NOW(), :creator_id
		)
	`, menu)
	if err != nil {
		return 0, fmt.Errorf("insert sys_config_menus failed: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %v", err)
	}
	menu.ID = uint(id)

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction failed: %v", err)
	}

	return menu.ID, nil
}

// UpdateMenu 更新菜单配置
func (s *MenuService) UpdateMenu(menu *model.ConfigMenu, updaterID uint) error {

	menu.UpdaterID = updaterID
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 更新菜单配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_menus SET 
			app_id = :app_id, 
			parent_id = :parent_id, 
			node_id = :node_id, 
			menu_name = :menu_name, 
			menu_code = :menu_code,
			menu_type = :menu_type,
			level = :level,
			sort = :sort,
			icon = :icon,
			path = :path,
			updated_at = NOW(),
			updater_id = :updater_id
		WHERE id = :id
	`, menu)
	if err != nil {
		return fmt.Errorf("update sys_config_menus failed: %v", err)
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
	err := s.db.Get(&menu, "SELECT * FROM sys_config_menus WHERE id = ?", id)
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
	_, err = tx.Exec("UPDATE sys_config_menus SET status = 0, updated_at = NOW() WHERE id = ?", id)
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
func (s *MenuService) ListMenus(appID uint, userID uint) ([]model.ConfigMenu, error) {
	var menus []model.ConfigMenu
	query := `
        SELECT m.* FROM sys_config_menus m
        INNER JOIN sys_permissions p ON m.id = p.menu_id
        INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
        INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ?
        AND m.app_id = ?
        AND p.status = 1
        ORDER BY m.sort ASC, m.id ASC
    `
	err := s.db.Select(&menus, query, userID, appID)
	if err != nil {
		return nil, fmt.Errorf("list menus failed: %v", err)
	}

	// 获取所有子菜单和父级菜单
	result := make(map[uint]*model.ConfigMenu)
	queryIDs := make(map[uint]struct{})
	for _, menu := range menus {
		result[menu.ID] = &menu
		if menu.ParentID != 0 {
			nodeIDs := strings.Split(menu.NodeID, "_")
			for _, nodeID := range nodeIDs {
				id := utils.ParseUint(nodeID)
				if _, ok := result[id]; !ok {
					queryIDs[id] = struct{}{}
				}
			}
		}
	}

	if len(queryIDs) > 0 {
		ids := make([]uint, 0, len(queryIDs))
		for id := range queryIDs {
			ids = append(ids, id)
		}

		query, args, err := sqlx.In(`
            SELECT * FROM sys_config_menus WHERE id IN (?)
        `, ids)
		if err != nil {
			return nil, err
		}

		query = s.db.Rebind(query)
		var parents []model.ConfigMenu
		err = s.db.Select(&parents, query, args...)
		if err != nil {
			return nil, err
		}
		for _, parent := range parents {
			result[parent.ID] = &parent
		}
	}

	for _, menu := range menus {
		query := `
            SELECT * FROM sys_config_menus WHERE node_id LIKE ? ORDER BY sort ASC, id ASC
        `
		var children []model.ConfigMenu
		err := s.db.Select(&children, query, menu.NodeID+"%")
		if err != nil {
			return nil, fmt.Errorf("list children menus failed: %v", err)
		}
		for _, child := range children {
			result[child.ID] = &child
		}
	}

	// 将结果转换为数组
	resultList := make([]model.ConfigMenu, 0, len(result))
	for _, menu := range result {
		resultList = append(resultList, *menu)
	}

	// 按照 id 字段排序
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].ID < resultList[j].ID
	})

	return resultList, nil
}

// TreeMenus 获取菜单树形结构
func (s *MenuService) TreeMenus(appID uint, userID uint) ([]model.TreeConfigMenu, error) {
	// 调用 ListMenus 获取所有相关的菜单
	menus, err := s.ListMenus(appID, userID)
	if err != nil {
		return nil, fmt.Errorf("list menus failed: %v", err)
	}

	// 创建一个map来存储所有菜单
	menuMap := make(map[uint]*model.TreeConfigMenu)
	for _, menu := range menus {
		menuMap[menu.ID] = &model.TreeConfigMenu{
			ID:        menu.ID,
			AppID:     menu.AppID,
			NodeID:    menu.NodeID,
			ParentID:  menu.ParentID,
			MenuName:  menu.MenuName,
			MenuCode:  menu.MenuCode,
			MenuType:  menu.MenuType,
			Level:     menu.Level,
			Sort:      menu.Sort,
			Icon:      menu.Icon,
			Path:      menu.Path,
			Status:    menu.Status,
			CreatedAt: menu.CreatedAt,
			CreatorID: menu.CreatorID,
			UpdatedAt: menu.UpdatedAt,
			UpdaterID: menu.UpdaterID,
			Children:  []*model.TreeConfigMenu{},
		}
	}

	// 创建树形结构
	var treeMenus []model.TreeConfigMenu
	for _, menu := range menuMap {
		if menu.ParentID == 0 {
			treeMenus = append(treeMenus, *menu)
		} else {
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, menu)
			}
		}
	}

	// 按照 id 字段排序
	sort.Slice(treeMenus, func(i, j int) bool {
		return treeMenus[i].ID < treeMenus[j].ID
	})

	return treeMenus, nil
}
