package rbac

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// UserRoleService 用户角色关联服务
type UserRoleService struct{}

// GetUserRoles 获取用户的角色列表
func (s *UserRoleService) GetUserRoles(userID uint) ([]map[string]interface{}, error) {
	var roles []map[string]interface{}
	query := `
		SELECT r.*
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ?
		ORDER BY r.created_at DESC
	`
	err := model.DB.Select(&roles, query, userID)
	return roles, err
}

// AddRoleToUser 为用户添加角色
func (s *UserRoleService) AddRoleToUser(operatorID uint, userID uint, roleID []uint) error {
	// 检查用户是否存在
	var userCount int
	err := model.DB.Get(&userCount, "SELECT COUNT(*) FROM sys_users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	if userCount == 0 {
		return errors.New("用户不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 添加新角色
	for _, rID := range roleID {
		// 检查角色是否存在
		var roleCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM sys_roles WHERE id = ?", rID).Scan(&roleCount)
		if err != nil {
			return err
		}
		if roleCount == 0 {
			return errors.New("角色不存在")
		}

		// 检查用户是否已拥有该角色
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND role_id = ?", userID, rID).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue // 用户已拥有该角色
		}

		_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id, creator_id, created_at) VALUES (?, ?, ?, ?)", userID, rID, operatorID, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// RemoveRolesFromUser 从用户移除角色
func (s *UserRoleService) RemoveRolesFromUser(operatorID, userID uint, roleID []uint) error {
	if len(roleID) == 0 {
		return errors.New("角色代码列表为空")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 移除角色
	for _, rID := range roleID {
		_, err = tx.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, rID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
