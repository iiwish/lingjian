package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/store"
	"github.com/iiwish/lingjian/pkg/utils"
	"github.com/mojocn/base64Captcha"
)

type AuthService struct {
	store store.Store
}

// AuthorizeRequest OAuth2授权请求参数
type AuthorizeRequest struct {
	// 客户端ID
	ClientID string `json:"client_id" binding:"required"`
	// 重定向URI
	RedirectURI string `json:"redirect_uri" binding:"required"`
	// 响应类型
	ResponseType string `json:"response_type" binding:"required"`
	// 权限范围
	Scope string `json:"scope" binding:"required"`
	// 状态
	State string `json:"state"`
	// 是否同意授权
	Approved bool `json:"approved"`
}

// TokenRequest OAuth2令牌请求参数
type TokenRequest struct {
	// 授权类型
	GrantType string `json:"grant_type" binding:"required"`
	// 客户端ID
	ClientID string `json:"client_id" binding:"required"`
	// 客户端密钥
	ClientSecret string `json:"client_secret" binding:"required"`
	// 授权码
	Code string `json:"code"`
	// 重定向URI
	RedirectURI string `json:"redirect_uri"`
	// 刷新令牌
	RefreshToken string `json:"refresh_token"`
}

// SwitchRoleRequest 切换角色请求参数
type SwitchRoleRequest struct {
	// 角色代码
	RoleCode string `json:"role_code" binding:"required"`
}

func NewAuthService(s store.Store) *AuthService {
	return &AuthService{
		store: s,
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
	// 令牌类型
	TokenType string `json:"token_type" example:"Bearer"`
}

// ChangePasswordRequest 修改密码请求参数
type ChangePasswordRequest struct {
	// 旧密码
	OldPassword string `json:"old_password" binding:"required" example:"123456"`
	// 新密码
	NewPassword string `json:"new_password" binding:"required" example:"1234567"`
}

// GenerateCaptcha 生成验证码
func (s *AuthService) GenerateCaptcha() (*CaptchaResponse, error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(40, 120, 4, 0.3, 50)
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
	// 打印请求参数
	log.Printf("LoginRequest: %+v", req)

	var user model.User
	err := model.DB.Get(&user, "SELECT * FROM sys_users WHERE username = ?", req.Username)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	if utils.HashPassword(req.Password) != user.Password {
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
	// 验证刷新令牌并获取claims
	_, err := utils.ParseToken(refreshToken, utils.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 验证令牌是否在存储中
	userId, err := s.store.VerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	var user model.User
	err = model.DB.Get(&user, "SELECT * FROM sys_users WHERE id = ?", userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成包含角色信息的新令牌
	tokenClaims := map[string]interface{}{
		"user_id":  userId,
		"username": user.Username,
	}

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateTokenWithClaims(tokenClaims, utils.AccessToken)
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

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userId uint, req *ChangePasswordRequest) error {
	var user model.User
	err := model.DB.Get(&user, "SELECT * FROM sys_users WHERE id = ?", userId)
	if err != nil {
		return errors.New("用户不存在")
	}

	if utils.HashPassword(req.OldPassword) != user.Password {
		return errors.New("旧密码错误")
	}

	// 更新密码
	_, err = model.DB.Exec("UPDATE sys_users SET password = ? WHERE id = ?", utils.HashPassword(req.NewPassword), userId)
	if err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}

// Logout 用户登出
func (s *AuthService) Logout(userId uint) error {
	return s.store.RemoveUserTokens(userId)
}

// HandleAuthorize 处理OAuth2授权请求
func (s *AuthService) HandleAuthorize(req *AuthorizeRequest) (string, error) {
	// 验证客户端
	if err := s.validateClient(req.ClientID, req.RedirectURI); err != nil {
		return "", err
	}

	// 如果用户拒绝授权
	if !req.Approved {
		return fmt.Sprintf("%s?error=access_denied&state=%s",
			req.RedirectURI, req.State), nil
	}

	// 生成授权码
	code := utils.GenerateRandomString(32)

	// 存储授权码（10分钟有效期）
	if err := s.store.StoreAuthCode(code, req.ClientID, req.Scope, 600); err != nil {
		return "", errors.New("存储授权码失败")
	}

	// 构建重定向URL
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s",
		req.RedirectURI, code, req.State)

	return redirectURL, nil
}

// HandleToken 处理OAuth2令牌请求
func (s *AuthService) HandleToken(req *TokenRequest) (*TokenResponse, error) {
	// 验证客户端
	if err := s.validateClientCredentials(req.ClientID, req.ClientSecret); err != nil {
		return nil, err
	}

	switch req.GrantType {
	case "authorization_code":
		return s.handleAuthorizationCode(req)
	case "refresh_token":
		return s.handleRefreshTokenGrant(req)
	default:
		return nil, errors.New("不支持的授权类型")
	}
}

// validateClient 验证客户端
func (s *AuthService) validateClient(clientID, redirectURI string) error {
	// TODO: 从数据库验证客户端信息
	if clientID != "test_client" || redirectURI != "http://localhost:3000/callback" {
		return errors.New("无效的客户端")
	}
	return nil
}

// validateClientCredentials 验证客户端凭证
func (s *AuthService) validateClientCredentials(clientID, clientSecret string) error {
	// TODO: 从数据库验证客户端凭证
	if clientID != "test_client" || clientSecret != "test_secret" {
		return errors.New("无效的客户端凭证")
	}
	return nil
}

// handleAuthorizationCode 处理授权码方式
func (s *AuthService) handleAuthorizationCode(req *TokenRequest) (*TokenResponse, error) {
	// 验证授权码
	clientID, scope, err := s.store.GetAuthCode(req.Code)
	if err != nil {
		return nil, errors.New("无效的授权码")
	}
	if clientID != req.ClientID {
		return nil, errors.New("授权码不匹配")
	}

	// 生成访问令牌和刷新令牌
	accessToken := utils.GenerateRandomString(32)
	refreshToken := utils.GenerateRandomString(32)

	// 存储令牌
	if err := s.store.StoreOAuthToken(accessToken, refreshToken, clientID, scope); err != nil {
		return nil, errors.New("存储令牌失败")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	}, nil
}

// handleRefreshTokenGrant 处理刷新令牌方式
func (s *AuthService) handleRefreshTokenGrant(req *TokenRequest) (*TokenResponse, error) {
	// 验证刷新令牌
	clientID, _, err := s.store.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}
	if clientID != req.ClientID {
		return nil, errors.New("刷新令牌不匹配")
	}

	// 生成新的访问令牌
	newAccessToken := utils.GenerateRandomString(32)

	// 存储新的访问令牌
	if err := s.store.UpdateOAuthAccessToken(req.RefreshToken, newAccessToken); err != nil {
		return nil, errors.New("更新访问令牌失败")
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	}, nil
}
