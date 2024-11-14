package rbac

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// RolePermissionService 角色权限关联服务
type RolePermissionService struct{}

// GetRolePermissions 获取角色的权限列表
func (s *RolePermissionService) GetRolePermissions(roleID string) ([]model.Permission, error) {
	if roleID == "" {
		return nil, errors.New("角色代码不能为空")
	}

	var permissions []model.Permission
	query := `
		SELECT p.*
		FROM sys_permissions p
		JOIN sys_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.created_at DESC
	`
	err := model.DB.Select(&permissions, query, roleID)
	return permissions, err
}

// AddPermissionsToRole 为角色添加权限
func (s *RolePermissionService) AddPermissionsToRole(operatorID uint, roleID string, permissions []uint) error {
	if len(permissions) == 0 {
		return errors.New("权限代码列表为空")
	}

	// 检查角色是否存在
	var roleCount int
	err := model.DB.Get(&roleCount, "SELECT COUNT(*) FROM sys_roles WHERE code = ?", roleID)
	if err != nil {
		return err
	}
	if roleCount == 0 {
		return errors.New("角色不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 添加新权限
	for _, permID := range permissions {
		// 检查权限是否存在
		var permCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM sys_permissions WHERE id = ?", permID).Scan(&permCount)
		if err != nil {
			return err
		}
		if permCount == 0 {
			return errors.New("权限不存在")
		}

		// 检查是否已经分配
		err = tx.QueryRow("SELECT COUNT(*) FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permID).Scan(&permCount)
		if err != nil {
			return err
		}
		if permCount > 0 {
			continue // 已存在的权限跳过
		}

		_, err = tx.Exec(`
			INSERT INTO role_permissions (role_id, permission_id, creator_id, created_at)
			VALUES (?, ?, ?, ?)
		`, roleID, permID, operatorID, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// RemovePermissionsFromRole 从角色移除权限
func (s *RolePermissionService) RemovePermissionsFromRole(userID uint, roleID string, permissions []uint) error {
	if len(permissions) == 0 {
		return errors.New("权限代码列表为空")
	}

	// 检查角色是否存在
	var roleCount int
	err := model.DB.Get(&roleCount, "SELECT COUNT(*) FROM sys_roles WHERE code = ?", roleID)
	if err != nil {
		return err
	}
	if roleCount == 0 {
		return errors.New("角色不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 移除权限
	for _, permID := range permissions {
		_, err = tx.Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
