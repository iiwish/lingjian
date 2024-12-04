package rbac

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
)

// RoleService 角色服务
type RoleService struct{}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles() ([]model.Role, error) {
	var roles []model.Role
	err := model.DB.Select(&roles, "SELECT * FROM sys_roles ORDER BY created_at DESC")
	return roles, err
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(operatorID uint, req *model.Role) error {
	// 检查角色代码是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_roles WHERE code = ?", req.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色代码已存在")
	}

	// 检查父角色是否存在
	// var parentID *uint
	// if req.ParentID != 0 {
	// 	var parent struct {
	// 		ID uint `db:"id"`
	// 	}
	// 	err = model.DB.Get(&parent, "SELECT id FROM sys_roles WHERE id = ?", req.ParentID)
	// 	if err != nil {
	// 		return errors.New("父角色不存在")
	// 	}
	// 	parentID = &parent.ID
	// }

	// 创建角色
	_, err = model.DB.Exec(
		"INSERT INTO sys_roles (name, code, parent_id, description, status, created_at, creator_id, updated_at, updater_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Name, req.Code, req.ParentID, req.Description, 1, time.Now(), operatorID, time.Now(), operatorID,
	)

	return err
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(operatorID uint, roleID uint, req *model.Role) error {
	// 检查角色是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_roles WHERE id = ?", roleID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("角色不存在")
	}

	// 检查角色代码是否重复（排除自身）
	err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_roles WHERE code = ? AND id != ?", req.Code, roleID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色代码已存在")
	}

	// 更新角色
	_, err = model.DB.Exec(`
		UPDATE sys_roles 
		SET name = ?, code = ?, description = ?, updated_at = ?, updater_id = ?
		WHERE id = ?
	`, req.Name, req.Code, req.Description, time.Now(), operatorID, roleID)

	return err
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(operatorID uint, roleID uint) error {
	// 检查角色是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_roles WHERE id = ?", roleID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("角色不存在")
	}

	// 检查是否有子角色
	err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_roles WHERE parent_id = ?", roleID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子角色")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除角色的权限关联
	_, err = tx.Exec("DELETE FROM sys_role_permissions WHERE role_id = ?", roleID)
	if err != nil {
		return err
	}

	// 删除用户与角色的关联
	_, err = tx.Exec("DELETE FROM sys_user_roles WHERE role_id = ?", roleID)
	if err != nil {
		return err
	}

	// 删除角色
	_, err = tx.Exec("DELETE FROM sys_roles WHERE id = ?", roleID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}
