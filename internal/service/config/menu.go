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
	// menu.SourceID = 0

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

	// 补充node_id、level、sort字段
	if menu.ParentID == 0 {
		menu.NodeID = strconv.Itoa(int(menu.ID))
		menu.Level = 1
		menu.Sort = 1
	} else {
		parentMenu, err := s.GetMenuByID(menu.ParentID)
		if err != nil {
			return 0, fmt.Errorf("get parent menu failed: %v", err)
		}
		menu.NodeID = parentMenu.NodeID + "_" + strconv.Itoa(int(menu.ID))
		menu.Level = parentMenu.Level + 1
		// 获取同级菜单的最大排序值
		var maxSort int
		err = tx.Get(&maxSort, `
			SELECT IFNULL(MAX(sort),0) FROM sys_config_menus
			WHERE parent_id = ? AND level = ?
		`, menu.ParentID, menu.Level)
		if err != nil {
			return 0, fmt.Errorf("get max sort failed: %v", err)
		}
		menu.Sort = maxSort + 1

	}
	// 更新菜单配置
	_, err = tx.NamedExec(`
		UPDATE sys_config_menus SET
			node_id = :node_id,
			level = :level,
			sort = :sort	
		WHERE id = :id
	`, menu)
	if err != nil {
		return 0, fmt.Errorf("update sys_config_menus failed: %v", err)
	}

	// 如果父节点不是0，但是元素类型是4，那么需要增加一个父节点是0的元素
	if menu.MenuType == 4 && menu.ParentID != 0 {
		tempMenu := &model.ConfigMenu{
			AppID:     menu.AppID,
			ParentID:  0,
			MenuName:  menu.MenuName,
			MenuCode:  menu.MenuCode,
			MenuType:  1,
			Level:     1,
			Sort:      1,
			Icon:      "folder",
			SourceID:  0,
			Status:    1,
			CreatorID: creatorID,
			UpdaterID: creatorID,
		}
		result, err := tx.NamedExec(`
				INSERT INTO sys_config_menus (
					app_id, parent_id, node_id, menu_name, menu_code, menu_type, level, sort, icon, source_id, status, created_at, creator_id, updated_at, updater_id
				) VALUES (
					:app_id, :parent_id, :node_id, :menu_name, :menu_code, :menu_type, :level, :sort, :icon, :source_id, :status, NOW(), :creator_id, NOW(), :creator_id
				)
			`, tempMenu)
		if err != nil {
			return 0, fmt.Errorf("insert sys_config_menus failed: %v", err)
		}
		// 获取插入的ID
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("get last insert id failed: %v", err)
		}
		tempMenu.ID = uint(id)
		tempMenu.NodeID = strconv.Itoa(int(tempMenu.ID))
		_, err = tx.NamedExec(`
				UPDATE sys_config_menus SET
					node_id = :node_id
				WHERE id = :id
			`, tempMenu)
		if err != nil {
			return 0, fmt.Errorf("update sys_config_menus failed: %v", err)
		}
		// 更新原节点的source_id
		_, err = tx.Exec("UPDATE sys_config_menus SET source_id = ? WHERE id = ?", tempMenu.ID, menu.ID)
		if err != nil {
			return 0, fmt.Errorf("update sys_config_menus failed: %v", err)
		}
	}

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
			menu_name = :menu_name, 
			menu_code = :menu_code,
			menu_type = :menu_type,
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

	// 获取菜单配置
	menu, err := s.GetMenuByID(id)
	if err != nil {
		return fmt.Errorf("get menu failed: %v", err)
	}

	// 删除菜单配置
	_, err = tx.Exec("DELETE FROM sys_config_menus WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete menu failed: %v", err)
	}

	// 删除菜单对应的元素
	tableMap := map[int]string{
		1: "sys_config_menus",
		2: "sys_config_tables",
		3: "sys_config_dimensions",
		4: "sys_config_menus",
		5: "sys_config_models",
		6: "sys_config_forms",
	}

	if menu.MenuType == 2 || menu.MenuType == 3 {
		// 查询tablename
		var tableName string
		err = tx.Get(&tableName, "SELECT table_name FROM "+tableMap[menu.MenuType]+" WHERE id = ?", menu.SourceID)
		if err != nil {
			return fmt.Errorf("get table name failed: %v", err)
		}
		_, err = tx.Exec("DROP TABLE IF EXISTS " + tableName)
		if err != nil {
			return fmt.Errorf("drop table failed: %v", err)
		}
	}

	if table, ok := tableMap[menu.MenuType]; ok {
		var query string
		if menu.MenuType == 1 {
			query = "DELETE FROM " + table + " WHERE node_id LIKE ?"
			_, err = tx.Exec(query, menu.NodeID+"_%")
		} else if menu.MenuType == 4 {
			query = "DELETE FROM " + table + " WHERE node_id = '" + strconv.Itoa(int(menu.SourceID)) + "' OR node_id LIKE ?"
			fmt.Println(query)
			_, err = tx.Exec(query, strconv.Itoa(int(menu.SourceID))+"_%")
		} else {
			query = "DELETE FROM " + table + " WHERE id = ?"
			_, err = tx.Exec(query, menu.SourceID)
		}
		if err != nil {
			return fmt.Errorf("delete failed: %v", err)
		}
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

	// 递归设置Path字段
	var setPath func(menu *model.TreeConfigMenu, parentPath string)
	setPath = func(menu *model.TreeConfigMenu, parentPath string) {
		if parentPath == "" {
			menu.Path = menu.MenuCode
		} else {
			menu.Path = parentPath + "/" + menu.MenuCode
		}
		// 对子节点进行排序
		sort.Slice(menu.Children, func(i, j int) bool {
			return menu.Children[i].Sort < menu.Children[j].Sort
		})
		for _, child := range menu.Children {
			setPath(child, menu.Path)
		}
	}

	for _, menu := range treeMenus {
		setPath(menu, "")
	}

	// 对树形结构进行排序
	sort.Slice(treeMenus, func(i, j int) bool {
		return treeMenus[i].Sort < treeMenus[j].Sort
	})

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

// UpdateMenuItemSort 更新菜单项排序
func (s *MenuService) UpdateMenuItemSort(userID uint, updaterID uint, menuID uint, parentID uint, sort int) error {
	// 开启事务
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %v", err)
	}
	defer tx.Rollback()

	// 获取当前节点的node_id、parent_id和sort
	var oldNode struct {
		NodeID   string `db:"node_id"`
		ParentID uint   `db:"parent_id"`
		Sort     int    `db:"sort"`
	}
	err = tx.Get(&oldNode, "SELECT node_id, parent_id, sort FROM sys_config_menus WHERE id = ?", menuID)
	if err != nil {
		return fmt.Errorf("get node_id and sort failed: %v", err)
	}

	// 检查父节点是否变化
	if parentID != oldNode.ParentID {
		// 获取新的node_id
		var newNodeID string
		if parentID != 0 {
			var parent struct {
				NodeID string `db:"node_id"`
			}
			err = tx.Get(&parent, "SELECT node_id FROM sys_config_menus WHERE id = ?", parentID)
			if err != nil {
				return fmt.Errorf("get parent node_id failed: %v", err)
			}
			newNodeID = parent.NodeID + "_" + fmt.Sprint(menuID)
		} else {
			newNodeID = fmt.Sprint(menuID)
		}
		// 更新node_id
		_, err = tx.Exec("UPDATE sys_config_menus SET node_id = ? WHERE id = ?", newNodeID, menuID)
		if err != nil {
			return fmt.Errorf("update node_id failed: %v", err)
		}
	}

	// 如果父节点变更，需要先从旧父节点移除，再插入到新父节点
	if parentID != oldNode.ParentID {
		// 获取新的node_id
		var newNodeID string
		if parentID != 0 {
			var parent struct {
				NodeID string `db:"node_id"`
			}
			err = tx.Get(&parent, "SELECT node_id FROM sys_config_menus WHERE id = ?", parentID)
			if err != nil {
				return fmt.Errorf("get parent node_id failed: %v", err)
			}
			newNodeID = parent.NodeID + "_" + fmt.Sprint(menuID)
		} else {
			newNodeID = fmt.Sprint(menuID)
		}

		// 1. 在旧父节点下移除该节点
		if oldNode.Sort != -1 {
			_, err = tx.Exec(`
                UPDATE sys_config_menus 
                SET sort = sort - 1 
                WHERE parent_id = ? 
                AND sort > ?
            `, oldNode.ParentID, oldNode.Sort)
			if err != nil {
				return err
			}
		}

		// 2. 为新父节点的sort腾位置
		_, err = tx.Exec(`
            UPDATE sys_config_menus
            SET sort = sort + 1
            WHERE parent_id = ?
            AND sort >= ?
        `, parentID, sort)
		if err != nil {
			return err
		}

		// 3. 更新该节点的父节点以及sort
		_, err = tx.Exec(`
            UPDATE sys_config_menus
            SET parent_id = ?,
                sort = ?,
                node_id = ?
            WHERE id = ?
        `, parentID, sort, newNodeID, menuID)
		if err != nil {
			return err
		}
	} else {
		// 如果只是同一父节点内sort值变动，按原逻辑处理
		if sort != oldNode.Sort {
			_, err = tx.Exec("UPDATE sys_config_menus SET sort = -1 WHERE id = ?", menuID)
			if err != nil {
				return err
			}
			if sort < oldNode.Sort {
				_, err = tx.Exec(`
                    UPDATE sys_config_menus 
                    SET sort = sort + 1 
                    WHERE parent_id = ? 
                    AND sort >= ? 
                    AND sort < ?
                `, parentID, sort, oldNode.Sort)
			} else {
				_, err = tx.Exec(`
                    UPDATE sys_config_menus 
                    SET sort = sort - 1 
                    WHERE parent_id = ? 
                    AND sort > ? 
                    AND sort <= ?
                `, parentID, oldNode.Sort, sort)
			}
			if err != nil {
				return err
			}
			_, err = tx.Exec("UPDATE sys_config_menus SET sort = ? WHERE id = ?", sort, menuID)
			if err != nil {
				return err
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %v", err)
	}

	return nil
}

// GetSystemMenuID 获取系统菜单ID
func (s *MenuService) GetSystemMenuID(appID uint) (uint, error) {
	var id uint
	err := s.db.Get(&id, "SELECT id FROM sys_config_menus WHERE menu_code = 'system' AND app_id = ?", appID)
	if err != nil {
		return 0, fmt.Errorf("get system menu id failed: %v", err)
	}
	return id, nil
}
