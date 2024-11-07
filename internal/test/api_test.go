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

// setupTestRouter 设置测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// TODO: 注册测试路由
	// v1.RegisterAuthRoutes(r.Group("/api/v1"))
	// v1.RegisterRBACRoutes(r.Group("/api/v1"))
	// v1.RegisterAppRoutes(r.Group("/api/v1"))

	return r
}

// TestAPIFlow 测试完整的API流程
func TestAPIFlow(t *testing.T) {
	router := setupTestRouter()
	var token string

	// 1. 测试用户登录
	t.Run("用户登录", func(t *testing.T) {
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

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		token = response["token"].(string)
	})

	// 2. 测试创建应用
	t.Run("创建应用", func(t *testing.T) {
		appData := map[string]interface{}{
			"name":        "测试应用",
			"code":        "test_app",
			"description": "这是一个测试应用",
		}
		jsonData, err := json.Marshal(appData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/apps", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 3. 测试创建角色
	t.Run("创建角色", func(t *testing.T) {
		roleData := map[string]interface{}{
			"name": "应用管理员",
			"code": "app_admin",
		}
		jsonData, err := json.Marshal(roleData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 4. 测试分配权限
	t.Run("分配权限", func(t *testing.T) {
		permData := map[string]interface{}{
			"permission_ids": []int{1, 2, 3},
		}
		jsonData, err := json.Marshal(permData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/roles/1/permissions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 5. 测试创建定时任务
	t.Run("创建定时任务", func(t *testing.T) {
		taskData := map[string]interface{}{
			"name":        "数据清理任务",
			"cron":        "0 0 * * *",
			"sql":         "DELETE FROM logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY)",
			"description": "每天零点清理7天前的日志",
		}
		jsonData, err := json.Marshal(taskData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 6. 测试获取应用列表
	t.Run("获取应用列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/apps", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "data")
	})
}
