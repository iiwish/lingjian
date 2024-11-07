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

// TestHelper 测试辅助结构体
type TestHelper struct {
	Router *gin.Engine
	Token  string
}

// NewTestHelper 创建测试辅助对象
func NewTestHelper(t *testing.T) *TestHelper {
	helper := &TestHelper{
		Router: setupTestRouter(),
	}
	helper.login(t)
	return helper
}

// login 登录并获取token
func (h *TestHelper) login(t *testing.T) {
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	jsonData, err := json.Marshal(loginData)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	h.Token = response["token"].(string)
}

// MakeRequest 发送HTTP请求
func (h *TestHelper) MakeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonData, err := json.Marshal(body)
		assert.NoError(t, err)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
	}

	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)
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

// CreateTestData 创建测试数据
func (h *TestHelper) CreateTestData(t *testing.T) {
	// 创建测试角色
	roleData := map[string]string{
		"name": "测试角色",
		"code": "test_role",
	}
	w := h.MakeRequest(t, "POST", "/api/v1/roles", roleData)
	h.AssertSuccess(t, w)

	// 创建测试应用
	appData := map[string]interface{}{
		"name":        "测试应用",
		"code":        "test_app",
		"description": "这是一个测试应用",
	}
	w = h.MakeRequest(t, "POST", "/api/v1/apps", appData)
	h.AssertSuccess(t, w)

	// 创建测试权限
	permData := map[string]interface{}{
		"name":   "测试权限",
		"code":   "test_permission",
		"type":   "api",
		"path":   "/api/v1/test",
		"method": "GET",
	}
	w = h.MakeRequest(t, "POST", "/api/v1/permissions", permData)
	h.AssertSuccess(t, w)
}

// CleanTestData 清理测试数据
func (h *TestHelper) CleanTestData(t *testing.T) {
	// 清理测试角色
	w := h.MakeRequest(t, "DELETE", "/api/v1/roles/test_role", nil)
	h.AssertSuccess(t, w)

	// 清理测试应用
	w = h.MakeRequest(t, "DELETE", "/api/v1/apps/test_app", nil)
	h.AssertSuccess(t, w)

	// 清理测试权限
	w = h.MakeRequest(t, "DELETE", "/api/v1/permissions/test_permission", nil)
	h.AssertSuccess(t, w)
}
