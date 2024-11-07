package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		username   string
		password   string
		wantStatus int
	}{
		{
			name:       "正常登录",
			username:   "admin",
			password:   "admin123",
			wantStatus: http.StatusOK,
		},
		{
			name:       "用户名错误",
			username:   "wronguser",
			password:   "admin123",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "密码错误",
			username:   "admin",
			password:   "wrongpass",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构造登录请求数据
			loginData := map[string]string{
				"username": tt.username,
				"password": tt.password,
			}
			jsonData, err := json.Marshal(loginData)
			assert.NoError(t, err)

			// 创建测试请求
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			w := httptest.NewRecorder()

			// TODO: 设置路由并处理请求
			// router := setupTestRouter()
			// router.ServeHTTP(w, req)

			// 检查响应状态码
			assert.Equal(t, tt.wantStatus, w.Code)

			// 如果是成功登录，检查返回的token
			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "refresh_token")
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		refreshToken string
		wantStatus   int
		wantNewToken bool
	}{
		{
			name:         "正常刷新token",
			refreshToken: "valid_refresh_token",
			wantStatus:   http.StatusOK,
			wantNewToken: true,
		},
		{
			name:         "无效的refresh token",
			refreshToken: "invalid_refresh_token",
			wantStatus:   http.StatusUnauthorized,
			wantNewToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求
			req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
			req.Header.Set("Authorization", "Bearer "+tt.refreshToken)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// TODO: 设置路由并处理请求
			// router := setupTestRouter()
			// router.ServeHTTP(w, req)

			// 检查响应状态码
			assert.Equal(t, tt.wantStatus, w.Code)

			// 如果期望获得新token，检查响应内容
			if tt.wantNewToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}
		})
	}
}
