package rbac

import (
	"errors"

	"github.com/iiwish/lingjian/internal/model"
)

// PermissionService 权限服务
type PermissionService struct{}

// ListPermissions 获取权限列表
func (s *PermissionService) ListPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := model.DB.Select(&permissions, "SELECT * FROM sys_permissions ORDER BY created_at DESC")
	return permissions, err
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(operatorID uint, req *model.Permission) error {
	req.CreatorID = operatorID
	req.UpdaterID = operatorID

	// 检查权限代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_permissions WHERE code = ?", req.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("权限代码已存在")
	}

	// 创建权限
	_, err = model.DB.Exec(`
		INSERT INTO sys_permissions (name, code, type, path, method, menu_id, status, description, created_at, creator_id, updated_at, updater_id)
		VALUES (:name, :code, :type, :path, :method, :menu_id, :status, :description, NOW(), :creator_id, NOW(), :updater_id)
	`, req)

	return err
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(operatorID uint, permissionID uint, req *model.Permission) error {
	req.UpdaterID = operatorID
	// 检查权限是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_permissions WHERE id = ?", permissionID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("权限不存在")
	}

	// 检查权限代码是否重复（排除自身）
	err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_permissions WHERE code = ? AND id != ?", req.Code, permissionID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("权限代码已存在")
	}

	// 更新权限
	_, err = model.DB.Exec(`
		UPDATE sys_permissions 
		SET name = :name, 
		code = :code, 
		type = :type, 
		path = :path, 
		method = :method, 
		menu_id = :menu_id,
		description = :description, 
		updated_at = NOW(), 
		updater_id = :updater_id
		WHERE id = ?
	`, req)

	return err
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(operatorID uint, permissionID uint) error {
	// 检查权限是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_permissions WHERE id = ?", permissionID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("权限不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除角色权限关联
	_, err = tx.Exec("DELETE FROM role_permissions WHERE permission_id = ?", permissionID)
	if err != nil {
		return err
	}

	// 删除权限
	_, err = tx.Exec("DELETE FROM sys_permissions WHERE id = ?", permissionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
