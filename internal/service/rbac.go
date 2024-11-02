package service

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

type RBACService struct{}

// CreateRole 创建角色
func (s *RBACService) CreateRole(name, code string) error {
	_, err := model.DB.Exec(`
		INSERT INTO roles (name, code, status, created_at, updated_at)
		VALUES (?, ?, 1, ?, ?)
	`, name, code, time.Now(), time.Now())
	return err
}

// CreatePermission 创建权限
func (s *RBACService) CreatePermission(name, code, typ, path, method string) error {
	_, err := model.DB.Exec(`
		INSERT INTO permissions (name, code, type, path, method, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?)
	`, name, code, typ, path, method, time.Now(), time.Now())
	return err
}

// AssignRoleToUser 为用户分配角色
func (s *RBACService) AssignRoleToUser(userID, roleID uint) error {
	// 检查用户是否存在
	var userExists bool
	err := model.DB.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在
	var roleExists bool
	err = model.DB.Get(&roleExists, "SELECT EXISTS(SELECT 1 FROM roles WHERE id = ?)", roleID)
	if err != nil {
		return err
	}
	if !roleExists {
		return errors.New("角色不存在")
	}

	// 分配角色
	_, err = model.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role_id = role_id
	`, userID, roleID)
	return err
}

// AssignPermissionToRole 为角色分配权限
func (s *RBACService) AssignPermissionToRole(roleID, permissionID uint) error {
	// 检查角色是否存在
	var roleExists bool
	err := model.DB.Get(&roleExists, "SELECT EXISTS(SELECT 1 FROM roles WHERE id = ?)", roleID)
	if err != nil {
		return err
	}
	if !roleExists {
		return errors.New("角色不存在")
	}

	// 检查权限是否存在
	var permExists bool
	err = model.DB.Get(&permExists, "SELECT EXISTS(SELECT 1 FROM permissions WHERE id = ?)", permissionID)
	if err != nil {
		return err
	}
	if !permExists {
		return errors.New("权限不存在")
	}

	// 分配权限
	_, err = model.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE permission_id = permission_id
	`, roleID, permissionID)
	return err
}

// GetUserRoles 获取用户的所有角色
func (s *RBACService) GetUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := model.DB.Select(&roles, `
		SELECT r.* FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.status = 1
	`, userID)
	return roles, err
}

// GetRolePermissions 获取角色的所有权限
func (s *RBACService) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := model.DB.Select(&permissions, `
		SELECT p.* FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ? AND p.status = 1
	`, roleID)
	return permissions, err
}
