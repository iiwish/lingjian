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
	// 基础查询,添加权限控制
	baseQuery := `
        SELECT DISTINCT m.* FROM sys_config_menus m
        INNER JOIN sys_permissions p ON m.id = p.menu_id
        INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
        INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ? 
        AND m.app_id = ?
        AND m.status = 1
        AND p.status = 1
    `
	args := []interface{}{userID, appID}

	// 添加level过滤
	if level != nil {
		baseQuery += " AND m.level = ?"
		args = append(args, *level)
	}

	// 添加parentID和menuType过滤
	if parentID != nil {
		if menuType == "descendants" {
			baseQuery += " AND (m.node_id LIKE ? OR m.node_id LIKE ?)"
			args = append(args, fmt.Sprintf("%d_%%", *parentID), fmt.Sprintf("%%_%d_%%", *parentID))
		} else {
			baseQuery += " AND m.parent_id = ?"
			args = append(args, *parentID)
		}
	}

	baseQuery += " ORDER BY m.sort ASC, m.id ASC"

	var menus []model.TreeConfigMenu
	err := s.db.Select(&menus, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("get menus failed: %v", err)
	}

	// 如果只需要指定level的菜单,直接返回
	if level != nil {
		return menus, nil
	}

	// 构建树形结构
	menuMap := make(map[uint]*model.TreeConfigMenu)
	for _, menu := range menus {
		menuCopy := menu // 创建副本避免指针问题
		menuMap[menu.ID] = &menuCopy
		menuMap[menu.ID].Children = []*model.TreeConfigMenu{}
	}

	// 如果是查询子集,需要额外查询父级菜单以构建完整的树形结构
	if parentID != nil && len(menus) > 0 {
		parentQuery := `
            SELECT DISTINCT m.* FROM sys_config_menus m
            INNER JOIN sys_permissions p ON m.id = p.menu_id
            INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
            INNER JOIN sys_user_roles ur ON rp.role_id = ur.role_id
            WHERE ur.user_id = ? 
            AND m.app_id = ?
            AND m.status = 1
            AND p.status = 1
            AND m.id = ?
        `
		var parent model.TreeConfigMenu
		err := s.db.Get(&parent, parentQuery, userID, appID, *parentID)
		if err == nil {
			menuMap[parent.ID] = &parent
			menuMap[parent.ID].Children = []*model.TreeConfigMenu{}
		}
	}

	// 构建树形结构
	var treeMenus []model.TreeConfigMenu
	if menuType == "children" && parentID != nil {
		// 只返回直接子集
		for _, menu := range menuMap {
			if menu.ID == *parentID {
				children := make([]model.TreeConfigMenu, 0)
				for _, child := range menu.Children {
					children = append(children, *child)
				}
				return children, nil
			}
		}
		return []model.TreeConfigMenu{}, nil
	} else {
		// 构建完整的树形结构
		for _, menu := range menuMap {
			if menu.ParentID == 0 || (parentID != nil && menu.ID == *parentID) {
				treeMenus = append(treeMenus, *menu)
			} else {
				if parent, ok := menuMap[menu.ParentID]; ok {
					parent.Children = append(parent.Children, menu)
				}
			}
		}

		// 如果是descendants类型,只返回指定parentID的子树
		if menuType == "descendants" && parentID != nil {
			for _, menu := range treeMenus {
				if menu.ID == *parentID {
					return []model.TreeConfigMenu{menu}, nil
				}
			}
			return []model.TreeConfigMenu{}, nil
		}
	}

	return treeMenus, nil
}
