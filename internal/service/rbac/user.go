package rbac

import (
	"errors"
	"time"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

// UserService 用户服务
type UserService struct{}

// GetUsers 获取用户列表
func (s *UserService) GetUsers() ([]model.User, error) {
	var users []model.User
	err := model.DB.Select(&users, "SELECT * FROM sys_users ORDER BY created_at DESC")
	return users, err
}

// GetUser 获取用户信息
func (s *UserService) GetUser(userID uint) (*model.User, error) {
	var user model.User
	err := model.DB.Get(&user, "SELECT * FROM sys_users WHERE id = ?", userID)
	return &user, err
}

// CreateUser 创建用户
func (s *UserService) CreateUser(operatorID uint, req *model.User) error {
	// 检查用户名是否已存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE username = ?", req.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE email = ?", req.Email)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已存在")
		}
	}
	// 检查手机号是否已存在
	if req.Phone != "" {
		err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE phone = ?", req.Phone)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}

	// 密码加密
	req.Password = utils.HashPassword(req.Password)

	// 创建用户
	_, err = model.DB.Exec(`
		INSERT INTO sys_users (username, password, nickname, email, phone, status, created_at, creator_id, updated_at, updater_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Username, req.Password, req.Nickname, req.Email, req.Phone, 1, time.Now(), operatorID, time.Now(), operatorID)

	return err
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(operatorID uint, userID uint, req *model.User) error {
	// 检查用户是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("用户不存在")
	}

	// 检查用户名是否重复（排除自身）
	err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE username = ? AND id != ?", req.Username, userID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否重复（排除自身）
	if req.Email != "" {
		err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE email = ? AND id != ?", req.Email, userID)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已存在")
		}
	}

	// 检查手机号是否重复（排除自身）
	if req.Phone != "" {
		err = model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE phone = ? AND id != ?", req.Phone, userID)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}

	// 密码加密
	if req.Password != "" {
		req.Password = utils.HashPassword(req.Password)
	}

	// 构建动态 SQL 查询
	query := `
        UPDATE sys_users 
        SET username = ?, nickname = ?, email = ?, phone = ?, updated_at = ?, updater_id = ?`
	args := []interface{}{req.Username, req.Nickname, req.Email, req.Phone, time.Now(), operatorID}

	if req.Password != "" {
		query += `, password = ?`
		args = append(args, req.Password)
	}

	query += ` WHERE id = ?`
	args = append(args, userID)

	// 更新用户
	_, err = model.DB.Exec(query, args...)

	return err
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(operatorID uint, userID uint) error {
	// 检查用户是否存在
	var count int
	err := model.DB.Get(&count, "SELECT COUNT(*) FROM sys_users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("用户不存在")
	}

	// 开始事务
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除用户角色关联
	_, err = tx.Exec("DELETE FROM sys_user_roles WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// 删除用户
	_, err = tx.Exec("DELETE FROM sys_users WHERE id = ?", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
