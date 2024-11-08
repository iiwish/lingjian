package service

import (
	"errors"
	"strings"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// RBACService RBAC服务
type RBACService struct{}

// CreateRole 创建角色
func (s *RBACService) CreateRole(req *model.CreateRoleRequest) error {
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
func (s *RBACService) CreatePermission(req *model.CreatePermissionRequest) error {
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

// AssignPermissionsToRole 为角色分配权限
func (s *RBACService) AssignPermissionsToRole(roleCode, appCode string, permissionCodes []string) error {
	if len(permissionCodes) == 0 {
		return errors.New("权限代码列表不能为空")
	}

	// 获取角色ID
	var roleID uint
	err := model.DB.Get(&roleID, "SELECT id FROM roles WHERE code = ? AND app_code = ?", roleCode, appCode)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 构建IN查询的占位符
	placeholders := strings.Repeat("?,", len(permissionCodes))
	placeholders = placeholders[:len(placeholders)-1] // 移除最后一个逗号

	// 获取权限IDs
	query := "SELECT id FROM permissions WHERE code IN (" + placeholders + ") AND app_code = ?"
	args := make([]interface{}, len(permissionCodes)+1)
	for i, code := range permissionCodes {
		args[i] = code
	}
	args[len(permissionCodes)] = appCode

	rows, err := model.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var permissionIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return err
		}
		permissionIDs = append(permissionIDs, id)
	}

	if len(permissionIDs) != len(permissionCodes) {
		return errors.New("部分权限代码不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除现有权限
	_, err = tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID)
	if err != nil {
		return err
	}

	// 分配新权限
	for _, permID := range permissionIDs {
		_, err = tx.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES (?, ?)
		`, roleID, permID)
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
func (s *RBACService) GetRolePermissions(roleCode, appCode string) ([]map[string]interface{}, error) {
	var permissions []map[string]interface{}
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
