package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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
			app_id, parent_id, node_id, menu_name, menu_code, menu_type, level, sort, icon, source_id, status, created_at, creator_id, updated_at, updater_id
		) VALUES (
			:app_id, :parent_id, :node_id, :menu_name, :menu_code, :menu_type, :level, :sort, :icon, :source_id, :status, NOW(), :creator_id, NOW(), :creator_id
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
			source_id = :source_id,
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
func (s *MenuService) GetMenuByID(id uint) (*model.ConfigMenu, error) {
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

	// 删除菜单配置
	_, err = tx.Exec("DELETE FROM sys_config_menus WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete menu failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetMenus 获取菜单列表
func (s *MenuService) GetMenus(appID uint, userID uint, level *int, parentID *uint, menuType string) ([]model.TreeConfigMenu, error) {
	// 第一步：获取用户有权限的菜单列表，包括node_id
	permissionMenusQuery := `
        SELECT DISTINCT m.id, m.node_id FROM sys_config_menus m
        INNER JOIN sys_permissions p ON m.id = p.menu_id
        INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
        INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ? 
        AND m.app_id = ?
        AND m.status = 1
        AND p.status = 1
    `
	permArgs := []interface{}{userID, appID}

	var permMenus []struct {
		ID     uint   `db:"id"`
		NodeID string `db:"node_id"`
	}
	err := s.db.Select(&permMenus, permissionMenusQuery, permArgs...)
	if err != nil {
		return nil, fmt.Errorf("获取用户权限菜单失败: %v", err)
	}

	// 第二步：根据node_id提取所有相关的菜单ID
	menuIDSet := make(map[uint]struct{})
	nodeIDSet := make(map[string]struct{})
	for _, permMenu := range permMenus {
		// 添加菜单自身ID
		menuIDSet[permMenu.ID] = struct{}{}
		nodeIDSet[permMenu.NodeID] = struct{}{}

		// 解析node_id，提取所有父级菜单ID
		idStrs := strings.Split(permMenu.NodeID, "_")
		for _, idStr := range idStrs {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				continue
			}
			menuIDSet[uint(id)] = struct{}{}
		}
	}

	// 第三步：查找子节点，添加到菜单ID集合中
	for nodeID := range nodeIDSet {
		// 查找以node_id为前缀的子节点
		childMenusQuery := `
            SELECT id FROM sys_config_menus
            WHERE node_id LIKE ?
            AND app_id = ?
            AND status = 1
        `
		childNodeIDPattern := nodeID + "_%"
		var childMenuIDs []uint
		err := s.db.Select(&childMenuIDs, childMenusQuery, childNodeIDPattern, appID)
		if err != nil {
			return nil, fmt.Errorf("获取子菜单失败: %v", err)
		}
		for _, id := range childMenuIDs {
			menuIDSet[id] = struct{}{}
		}
	}

	// 将菜单ID集合转换为切片
	menuIDs := make([]uint, 0, len(menuIDSet))
	for id := range menuIDSet {
		menuIDs = append(menuIDs, id)
	}

	// 如果需要按level过滤
	if level != nil {
		levelMenusQuery, args, err := sqlx.In(`
            SELECT * FROM sys_config_menus
            WHERE id IN (?)
            AND level = ?
            ORDER BY sort ASC, id ASC
        `, menuIDs, *level)
		if err != nil {
			return nil, fmt.Errorf("构建查询语句失败: %v", err)
		}
		levelMenusQuery = s.db.Rebind(levelMenusQuery)

		var menus []model.TreeConfigMenu
		err = s.db.Select(&menus, levelMenusQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("获取菜单失败: %v", err)
		}

		return menus, nil
	}

	// 第四步：查询所有相关的菜单记录
	menusQuery, args, err := sqlx.In(`
        SELECT * FROM sys_config_menus
        WHERE id IN (?)
        ORDER BY sort ASC, id ASC
    `, menuIDs)
	if err != nil {
		return nil, fmt.Errorf("构建查询语句失败: %v", err)
	}
	menusQuery = s.db.Rebind(menusQuery)

	var menus []model.TreeConfigMenu
	err = s.db.Select(&menus, menusQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("获取菜单失败: %v", err)
	}

	// 第五步：构建树形结构
	menuMap := make(map[uint]*model.TreeConfigMenu)
	for i := range menus {
		menu := &menus[i]
		menu.Children = []*model.TreeConfigMenu{}
		menuMap[menu.ID] = menu
	}

	var treeMenus []*model.TreeConfigMenu
	for _, menu := range menuMap {
		if menu.ParentID == 0 || menuMap[menu.ParentID] == nil {
			treeMenus = append(treeMenus, menu)
		} else {
			parent := menuMap[menu.ParentID]
			parent.Children = append(parent.Children, menu)
		}
	}

	// 对树形结构进行排序
	sort.Slice(treeMenus, func(i, j int) bool {
		return treeMenus[i].Sort < treeMenus[j].Sort
	})
	for _, menu := range treeMenus {
		sort.Slice(menu.Children, func(i, j int) bool {
			return menu.Children[i].Sort < menu.Children[j].Sort
		})
	}

	// 如果parentID不为nil，只返回对应的子树
	if parentID != nil {
		if rootMenu, ok := menuMap[*parentID]; ok {
			if menuType == "children" {
				children := make([]model.TreeConfigMenu, len(rootMenu.Children))
				for i, child := range rootMenu.Children {
					children[i] = *child
				}
				return children, nil
			} else if menuType == "descendants" {
				return []model.TreeConfigMenu{*rootMenu}, nil
			}
		}
		return []model.TreeConfigMenu{}, nil
	}

	// 将 []*model.TreeConfigMenu 转换为 []model.TreeConfigMenu
	result := make([]model.TreeConfigMenu, len(treeMenus))
	for i, menu := range treeMenus {
		result[i] = *menu
	}

	return result, nil
}
