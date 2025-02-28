package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/iiwish/lingjian/api/v1"
	"github.com/iiwish/lingjian/api/v1/config"
	"github.com/iiwish/lingjian/internal/middleware"
	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/pkg/queue"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var mockStore *MockStore

func init() {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 加载测试配置
	viper.SetConfigName("config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading test config file: %s", err)
	}

	// 初始化测试数据库连接
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// 初始化Mock存储
	mockStore = &MockStore{}

	// 设置认证服务和中间件使用Mock存储
	v1.InitAuthService(mockStore)
	middleware.SetStore(mockStore)

	// 初始化RabbitMQ连接
	if err := queue.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize test RabbitMQ: %v", err)
	}
}

// hashPassword 密码加密
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// cleanTestData 清理测试数据
func cleanTestData() error {
	tables := []string{
		"sys_role_permissions",
		"sys_user_roles",
		"sys_user_apps",
		"sys_element_triggers",
		"sys_task_logs",
		"sys_scheduled_tasks",
		"sys_config_menu",
		"sys_config_form",
		"sys_config_model",
		"sys_config_dimension",
		"sys_config_table",
		"sys_permissions",
		"sys_roles",
		"sys_apps",
		"sys_users",
	}

	for _, table := range tables {
		_, err := model.DB.Exec("DELETE FROM " + table)
		if err != nil {
			log.Printf("Error cleaning table %s: %v", table, err)
			return err
		}
	}

	return nil
}

// initTestData 初始化测试数据
func initTestData() error {
	log.Println("开始初始化测试数据...")

	// 先清理现有数据
	if err := cleanTestData(); err != nil {
		log.Printf("清理数据失败: %v", err)
		return err
	}

	now := time.Now()

	// 创建测试用户（使用加密后的密码）
	hashedPassword := hashPassword("admin1324")
	log.Printf("创建测试用户，密码哈希: %s", hashedPassword)
	_, err := model.DB.Exec(`
		INSERT INTO sys_users (id, username, password, email, phone, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, 1, "admin", hashedPassword, "admin@test.com", "13800138000", 1, now, now)
	if err != nil {
		log.Printf("创建测试用户失败: %v", err)
		return fmt.Errorf("failed to create test user: %v", err)
	}

	// 创建测试应用
	log.Println("创建测试应用...")
	_, err = model.DB.Exec(`
		INSERT INTO sys_apps (name, code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?, ?)
	`,
		"测试应用1", "test_app1", "用于测试的应用1", 1, now, now,
		"测试应用2", "test_app2", "用于测试的应用2", 1, now, now,
	)
	if err != nil {
		log.Printf("创建测试应用失败: %v", err)
		return fmt.Errorf("failed to create test apps: %v", err)
	}

	// 创建测试角色
	log.Println("创建测试角色...")
	_, err = model.DB.Exec(`
		INSERT INTO sys_roles (name, code, app_code, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "管理员", "app_admin", "test_app1", "系统管理员", 1, now, now)
	if err != nil {
		log.Printf("创建测试角色失败: %v", err)
		return fmt.Errorf("failed to create test role: %v", err)
	}

	// 创建测试权限
	log.Println("创建测试权限...")
	_, err = model.DB.Exec(`
		INSERT INTO sys_permissions (name, code, app_code, type, path, method, description, status, created_at, updated_at)
		VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		"查看应用", "view_apps", "test_app1", "api", "/api/v1/apps", "GET", "查看应用列表", 1, now, now,
		"创建应用", "create_app", "test_app1", "api", "/api/v1/apps", "POST", "创建新应用", 1, now, now,
		"创建角色", "create_role", "test_app1", "api", "/api/v1/roles", "POST", "创建新角色", 1, now, now,
		"分配权限", "assign_permission", "test_app1", "api", "/api/v1/roles/*/permissions", "POST", "为角色分配权限", 1, now, now,
		"创建任务", "create_task", "test_app1", "api", "/api/v1/tasks/scheduled", "POST", "创建定时任务", 1, now, now,
	)
	if err != nil {
		log.Printf("创建测试权限失败: %v", err)
		return fmt.Errorf("failed to create test permissions: %v", err)
	}

	// 分配角色给用户
	log.Println("分配角色给用户...")
	_, err = model.DB.Exec(`
		INSERT INTO sys_user_roles (user_id, role_id)
		SELECT u.id, r.id
		FROM sys_users u, sys_roles r
		WHERE u.username = 'admin' AND r.code = 'app_admin'
	`)
	if err != nil {
		log.Printf("分配角色给用户失败: %v", err)
		return fmt.Errorf("failed to assign role to user: %v", err)
	}

	// 分配权限给角色
	log.Println("分配权限给角色...")
	_, err = model.DB.Exec(`
		INSERT INTO sys_role_permissions (role_id, permission_id)
		SELECT r.id, p.id
		FROM sys_roles r
		CROSS JOIN sys_permissions p
		WHERE r.code = 'app_admin' AND r.app_code = p.app_code
	`)
	if err != nil {
		log.Printf("分配权限给角色失败: %v", err)
		return fmt.Errorf("failed to assign permissions to role: %v", err)
	}

	log.Println("测试数据初始化完成")
	return nil
}

// TestHelper 测试辅助结构体
type TestHelper struct {
	Router       *gin.Engine
	Token        string
	RefreshToken string
}

// NewTestHelper 创建测试辅助对象
func NewTestHelper(t *testing.T) *TestHelper {
	// 初始化测试数据
	err := initTestData()
	if err != nil {
		t.Fatalf("Failed to initialize test data: %v", err)
	}

	helper := &TestHelper{
		Router: setupTestRouter(),
	}
	helper.login(t)
	return helper
}

// setupTestRouter 设置测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// API路由
	api := r.Group("/api")
	{
		// v1版本API
		v1Group := api.Group("/v1")
		{
			// 注册认证相关路由（不需要认证）
			auth := v1Group.Group("/auth")
			{
				auth.GET("/captcha", v1.GetCaptcha)
				auth.POST("/login", v1.Login)
				auth.POST("/refresh", v1.RefreshToken)
			}

			// 需要认证的路由
			authorized := v1Group.Group("/")
			authorized.Use(middleware.AuthMiddleware())
			{
				// 认证相关路由（需要认证但不需要角色）
				auth := authorized.Group("/auth")
				{
					auth.POST("/logout", v1.Logout)
				}

				// 需要RBAC权限控制的路由
				rbacProtected := authorized.Group("/")
				rbacProtected.Use(middleware.RBACMiddleware())
				{
					// 注册RBAC相关路由
					v1.RegisterRBACRoutes(rbacProtected)
					// 注册应用相关路由
					v1.RegisterAppRoutes(rbacProtected)
					// 注册配置相关路由
					config.RegisterConfigRoutes(rbacProtected)
					// 注册任务相关路由
					v1.RegisterTaskRoutes(rbacProtected)
				}
			}
		}
	}

	return r
}

// login 登录并获取token
func (h *TestHelper) login(t *testing.T) {
	// 先获取验证码
	req := httptest.NewRequest("GET", "/api/v1/auth/captcha", nil)
	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	var captchaResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &captchaResp)
	assert.NoError(t, err, "Failed to parse captcha response")

	data := captchaResp["data"].(map[string]interface{})
	captchaId := data["captcha_id"].(string)

	// 登录请求
	loginData := map[string]interface{}{
		"username":    "admin",
		"password":    "admin123",
		"captcha_id":  captchaId,
		"captcha_val": "1234", // 在测试模式下，验证码固定为1234
	}
	jsonData, err := json.Marshal(loginData)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	// 打印响应内容以便调试
	t.Logf("Login Response Status: %d", w.Code)
	t.Logf("Login Response Body: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code, "Login request failed")

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to parse login response")

	// 检查响应结构
	data = response["data"].(map[string]interface{})
	assert.Contains(t, data, "access_token", "Response does not contain access_token")
	assert.Contains(t, data, "refresh_token", "Response does not contain refresh_token")

	token, ok := data["access_token"].(string)
	assert.True(t, ok, "Token is not a string")
	assert.NotEmpty(t, token, "Token is empty")

	refreshToken, ok := data["refresh_token"].(string)
	assert.True(t, ok, "Refresh token is not a string")
	assert.NotEmpty(t, refreshToken, "Refresh token is empty")

	h.Token = token
	h.RefreshToken = refreshToken

	// 打印token信息以便调试
	t.Logf("Access Token: %s", token)
	t.Logf("Refresh Token: %s", refreshToken)

	// 设置默认角色
	switchRoleData := map[string]interface{}{
		"role_code": "app_admin",
	}
	jsonData, err = json.Marshal(switchRoleData)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/api/v1/auth/switch-role", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Switch role request failed")

	// 解析响应
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to parse switch role response")

	// 检查响应结构
	if response["code"].(float64) == 200 {
		data = response["data"].(map[string]interface{})
		if data != nil {
			// 更新token为带有角色信息的新token
			h.Token = data["access_token"].(string)
			h.RefreshToken = data["refresh_token"].(string)
			t.Logf("New Access Token after role switch: %s", h.Token)
		}
	} else {
		t.Logf("Switch role failed with response: %s", w.Body.String())
	}
}

// MakeRequest 发送HTTP请求
func (h *TestHelper) MakeRequest(t *testing.T, method, path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonData, err := json.Marshal(body)
		assert.NoError(t, err)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	// 设置默认的Authorization头
	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
		// 打印请求头以便调试
		t.Logf("Request Authorization Header: %s", req.Header.Get("Authorization"))
	}

	// 设置额外的头部
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	// 打印请求和响应信息以便调试
	t.Logf("Request: %s %s", method, path)
	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())

	return w
}

// ParseResponse 解析响应数据
func (h *TestHelper) ParseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

// AssertSuccess 断言请求成功
func (h *TestHelper) AssertSuccess(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	assert.Equal(t, http.StatusOK, w.Code)
	return h.ParseResponse(t, w)
}

// AssertError 断言请求失败
func (h *TestHelper) AssertError(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, w.Code)
	return h.ParseResponse(t, w)
}
