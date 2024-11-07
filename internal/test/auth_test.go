package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)
	router := setupTestRouter()

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
		{
			name:       "用户名为空",
			username:   "",
			password:   "admin123",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "密码为空",
			username:   "admin",
			password:   "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginData := map[string]string{
				"username": tt.username,
				"password": tt.password,
			}
			jsonData, err := json.Marshal(loginData)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

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
	gin.SetMode(gin.TestMode)
	router := setupTestRouter()

	// 先登录获取有效的refresh token
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	jsonData, err := json.Marshal(loginData)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	validRefreshToken := loginResponse["refresh_token"].(string)

	tests := []struct {
		name         string
		refreshToken string
		wantStatus   int
		wantNewToken bool
	}{
		{
			name:         "正常刷新token",
			refreshToken: validRefreshToken,
			wantStatus:   http.StatusOK,
			wantNewToken: true,
		},
		{
			name:         "无效的refresh token",
			refreshToken: "invalid_refresh_token",
			wantStatus:   http.StatusUnauthorized,
			wantNewToken: false,
		},
		{
			name:         "refresh token为空",
			refreshToken: "",
			wantStatus:   http.StatusUnauthorized,
			wantNewToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
			if tt.refreshToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.refreshToken)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantNewToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}
		})
	}
}

func TestOAuth2Flow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupTestRouter()

	// 1. 测试获取授权页面
	t.Run("获取授权页面", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/auth/oauth/authorize?client_id=test_client&redirect_uri=http://localhost:3000/callback&response_type=code&scope=read", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "授权页面")
	})

	// 2. 测试授权确认
	var authCode string
	t.Run("授权确认", func(t *testing.T) {
		authData := map[string]interface{}{
			"client_id":     "test_client",
			"redirect_uri":  "http://localhost:3000/callback",
			"response_type": "code",
			"scope":         "read",
			"state":         "random_state",
			"approved":      true,
		}
		jsonData, err := json.Marshal(authData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/auth/oauth/authorize", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		location := w.Header().Get("Location")
		assert.Contains(t, location, "code=")
		assert.Contains(t, location, "state=random_state")

		// 从location中提取授权码
		// 简单处理，实际可能需要更复杂的URL解析
		authCode = location[strings.Index(location, "code=")+5 : strings.Index(location, "&state")]
	})

	// 3. 测试使用授权码获取token
	var accessToken string
	var refreshToken string
	t.Run("使用授权码获取token", func(t *testing.T) {
		tokenData := map[string]interface{}{
			"grant_type":    "authorization_code",
			"client_id":     "test_client",
			"client_secret": "test_secret",
			"code":          authCode,
			"redirect_uri":  "http://localhost:3000/callback",
		}
		jsonData, err := json.Marshal(tokenData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/auth/oauth/token", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
		assert.Contains(t, response, "expires_in")
		assert.Contains(t, response, "token_type")
		assert.Equal(t, "Bearer", response["token_type"])

		accessToken = response["access_token"].(string)
		refreshToken = response["refresh_token"].(string)
	})

	// 4. 测试使用refresh token刷新token
	t.Run("使用refresh token刷新token", func(t *testing.T) {
		refreshData := map[string]interface{}{
			"grant_type":    "refresh_token",
			"client_id":     "test_client",
			"client_secret": "test_secret",
			"refresh_token": refreshToken,
		}
		jsonData, err := json.Marshal(refreshData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/auth/oauth/token", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
		assert.Contains(t, response, "expires_in")
		assert.Contains(t, response, "token_type")
		assert.Equal(t, "Bearer", response["token_type"])
	})

	// 5. 测试使用access token访问受保护的资源
	t.Run("使用access token访问资源", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
