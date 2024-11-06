package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/store"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/mojocn/base64Captcha"
)

type AuthService struct {
	store *store.RedisStore
}

func NewAuthService() *AuthService {
	return &AuthService{
		store: store.NewRedisStore(),
	}
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	// 用户名
	Username string `json:"username" binding:"required" example:"admin"`
	// 密码
	Password string `json:"password" binding:"required" example:"123456"`
	// 验证码ID
	CaptchaId string `json:"captcha_id" binding:"required" example:"captcha-123"`
	// 验证码值
	CaptchaVal string `json:"captcha_val" binding:"required" example:"1234"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	// 访问令牌
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// 刷新令牌
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// 访问令牌过期时间（秒）
	ExpiresIn int `json:"expires_in" example:"7200"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	// 验证码ID
	CaptchaId string `json:"captcha_id" example:"captcha-123"`
	// Base64编码的验证码图片
	CaptchaImg string `json:"captcha_img" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	// 访问令牌
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// 刷新令牌
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// 访问令牌过期时间（秒）
	ExpiresIn int `json:"expires_in" example:"7200"`
}

// hashPassword 密码加密
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateCaptcha 生成验证码
func (s *AuthService) GenerateCaptcha() (*CaptchaResponse, error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(40, 120, 4, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, s.store)

	// 生成验证码
	id, b64s, _, err := c.Generate()
	if err != nil {
		return nil, errors.New("生成验证码失败")
	}

	return &CaptchaResponse{
		CaptchaId:  id,
		CaptchaImg: b64s,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 验证验证码
	if !s.store.Verify(req.CaptchaId, req.CaptchaVal, true) {
		return nil, errors.New("验证码错误")
	}

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

	// 存储令牌
	if err := s.store.StoreAccessToken(user.ID, accessToken); err != nil {
		return nil, err
	}
	if err := s.store.StoreRefreshToken(user.ID, refreshToken); err != nil {
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
	// 验证刷新令牌
	userId, err := s.store.VerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	var user model.User
	err = model.DB.Get(&user, "SELECT * FROM users WHERE id = ?", userId)
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

	// 存储新的访问令牌
	if err := s.store.StoreAccessToken(user.ID, newAccessToken); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userId uint) error {
	return s.store.RemoveUserTokens(userId)
}
