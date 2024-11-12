package service

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// RBACService RBAC服务
type RBACService struct{}

// CreateRole 创建角色
func (s *RBACService) CreateRole(req *model.Role) error {
	// 检查角色代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM roles WHERE code = ? AND app_code = ?", req.Code, req.AppCode)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色代码已存在")
	}

	// 检查父角色是否存在
	var parentID *uint
	if req.ParentCode != "" {
		var parent struct {
			ID      uint   `db:"id"`
			AppCode string `db:"app_code"`
		}
		err = model.DB.Get(&parent, "SELECT id, app_code FROM roles WHERE code = ?", req.ParentCode)
		if err != nil {
			return errors.New("父角色不存在")
		}
		if parent.AppCode != req.AppCode {
			return errors.New("父角色必须属于同一个应用")
		}
		parentID = &parent.ID
	}

	// 创建角色
	_, err = model.DB.Exec(`
		INSERT INTO roles (name, code, app_code, parent_id, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Name, req.Code, req.AppCode, parentID, req.Description, 1, time.Now(), time.Now())

	return err
}

// CreatePermission 创建权限
func (s *RBACService) CreatePermission(req *model.Permission) error {
	// 检查权限代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM permissions WHERE code = ? AND app_code = ?", req.Code, req.AppCode)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("权限代码已存在")
	}

	// 创建权限
	_, err = model.DB.Exec(`
		INSERT INTO permissions (name, code, app_code, type, path, method, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Name, req.Code, req.AppCode, req.Type, req.Path, req.Method, req.Description, 1, time.Now(), time.Now())

	return err
}

// AssignRoleToUser 为用户分配角色
func (s *RBACService) AssignRoleToUser(userID, roleID uint) error {
	// 检查用户和角色是否存在
	var userCount, roleCount int
	err := model.DB.Get(&userCount, "SELECT COUNT(*) FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	if userCount == 0 {
		return errors.New("用户不存在")
	}

	err = model.DB.Get(&roleCount, "SELECT COUNT(*) FROM roles WHERE id = ?", roleID)
	if err != nil {
		return err
	}
	if roleCount == 0 {
		return errors.New("角色不存在")
	}

	// 分配角色给用户
	_, err = model.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role_id = ?
	`, userID, roleID, roleID)

	return err
}

// AddPermissionsToRole 为角色添加权限
func (s *RBACService) AddPermissionsToRole(userID uint, roleID string, permissions []uint) error {
	if len(permissions) == 0 {
		return errors.New("权限代码列表为空")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 添加新权限
	for _, permID := range permissions {
		_, err = tx.Exec(`
			INSERT INTO role_permissions (role_id, permission_id, creator_id, created_at)
			VALUES (?, ?, ?, ?)
		`, roleID, permID, userID, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetUserRoles 获取用户的角色列表
func (s *RBACService) GetUserRoles(userID uint) ([]map[string]interface{}, error) {
	var roles []map[string]interface{}
	query := `
		SELECT r.*
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.status = 1
		ORDER BY r.created_at DESC
	`
	err := model.DB.Select(&roles, query, userID)
	return roles, err
}

// GetRolePermissions 获取角色的权限列表
func (s *RBACService) GetRolePermissions(roleCode, appCode string) ([]model.Permission, error) {
	var permissions []model.Permission
	query := `
		SELECT p.*
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON r.id = rp.role_id
		WHERE r.code = ? AND r.app_code = ? AND p.status = 1
		ORDER BY p.created_at DESC
	`
	err := model.DB.Select(&permissions, query, roleCode, appCode)
	return permissions, err
}
