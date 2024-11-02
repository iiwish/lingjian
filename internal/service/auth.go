package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/utils"
)

type AuthService struct{}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// hashPassword 密码加密
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.User
	err := model.DB.Get(&user, "SELECT * FROM users WHERE username = ?", req.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if hashPassword(req.Password) != user.Password {
		return nil, errors.New("密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成访问令牌
	accessToken, err := utils.GenerateToken(user.ID, user.Username, utils.AccessToken)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := utils.GenerateToken(user.ID, user.Username, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	claims, err := utils.ParseToken(refreshToken, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = model.DB.Get(&user, "SELECT * FROM users WHERE id = ?", claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateToken(user.ID, user.Username, utils.AccessToken)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
	}, nil
}
