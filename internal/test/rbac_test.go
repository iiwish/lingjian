package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBACFlow(t *testing.T) {
	helper := NewTestHelper(t)

	// 1. 创建多个应用
	apps := []struct {
		name string
		code string
	}{
		{"测试应用1", "test_app1"},
		{"测试应用2", "test_app2"},
	}

	for _, app := range apps {
		t.Run(fmt.Sprintf("创建应用-%s", app.name), func(t *testing.T) {
			appData := map[string]interface{}{
				"name":        app.name,
				"code":        app.code,
				"description": fmt.Sprintf("这是%s", app.name),
			}
			w := helper.MakeRequest(t, "POST", "/api/v1/apps", appData)
			helper.AssertSuccess(t, w)
		})
	}

	// 2. 创建角色层级关系
	roles := []struct {
		name     string
		code     string
		appCode  string
		parentID string
	}{
		{"超级管理员", "super_admin", "test_app1", ""},
		{"普通管理员", "admin", "test_app1", "super_admin"},
		{"普通用户", "user", "test_app1", "admin"},
		{"访客", "guest", "test_app1", "user"},
	}

	for _, role := range roles {
		t.Run(fmt.Sprintf("创建角色-%s", role.name), func(t *testing.T) {
			roleData := map[string]interface{}{
				"name":     role.name,
				"code":     role.code,
				"app_code": role.appCode,
			}
			if role.parentID != "" {
				roleData["parent_code"] = role.parentID
			}
			w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
			helper.AssertSuccess(t, w)
		})
	}

	// 3. 创建权限
	permissions := []struct {
		name        string
		code        string
		appCode     string
		type_       string
		path        string
		method      string
		description string
	}{
		{
			name:        "查看应用",
			code:        "view_apps",
			appCode:     "test_app1",
			type_:       "api",
			path:        "/api/v1/apps",
			method:      "GET",
			description: "允许查看应用列表",
		},
		{
			name:        "创建应用",
			code:        "create_app",
			appCode:     "test_app1",
			type_:       "api",
			path:        "/api/v1/apps",
			method:      "POST",
			description: "允许创建新应用",
		},
	}

	for _, perm := range permissions {
		t.Run(fmt.Sprintf("创建权限-%s", perm.name), func(t *testing.T) {
			permData := map[string]interface{}{
				"name":        perm.name,
				"code":        perm.code,
				"app_code":    perm.appCode,
				"type":        perm.type_,
				"path":        perm.path,
				"method":      perm.method,
				"description": perm.description,
			}
			w := helper.MakeRequest(t, "POST", "/api/v1/permissions", permData)
			helper.AssertSuccess(t, w)
		})
	}

	// 4. 分配权限给角色
	t.Run("分配权限给角色", func(t *testing.T) {
		// 超级管理员拥有所有权限
		permCodes := []string{"view_apps", "create_app"}
		assignData := map[string]interface{}{
			"permission_codes": permCodes,
			"app_code":         "test_app1",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles/super_admin/permissions", assignData)
		helper.AssertSuccess(t, w)

		// 普通管理员只有查看权限
		assignData = map[string]interface{}{
			"permission_codes": []string{"view_apps"},
			"app_code":         "test_app1",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/roles/admin/permissions", assignData)
		helper.AssertSuccess(t, w)
	})

	// 5. 测试权限继承
	t.Run("测试权限继承", func(t *testing.T) {
		// 检查超级管理员权限
		w := helper.MakeRequest(t, "GET", "/api/v1/roles/super_admin/permissions?app_code=test_app1", nil)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].(map[string]interface{})
		permissions := data["items"].([]interface{})
		assert.Equal(t, 2, len(permissions))

		// 检查普通管理员权限
		w = helper.MakeRequest(t, "GET", "/api/v1/roles/admin/permissions?app_code=test_app1", nil)
		resp = helper.AssertSuccess(t, w)
		data = resp["data"].(map[string]interface{})
		permissions = data["items"].([]interface{})
		assert.Equal(t, 1, len(permissions))
	})

	// 6. 测试权限检查
	t.Run("测试权限检查", func(t *testing.T) {
		// 超级管理员可以访问所有接口
		w := helper.MakeRequest(t, "GET", "/api/v1/apps", nil)
		helper.AssertSuccess(t, w)

		w = helper.MakeRequest(t, "POST", "/api/v1/apps", map[string]interface{}{
			"name":        "测试应用3",
			"code":        "test_app3",
			"description": "这是测试应用3",
		})
		helper.AssertSuccess(t, w)

		// 切换到普通管理员角色
		w = helper.MakeRequest(t, "POST", "/api/v1/auth/switch-role", map[string]string{"role_code": "admin"})
		helper.AssertSuccess(t, w)

		// 普通管理员只能查看
		w = helper.MakeRequest(t, "GET", "/api/v1/apps", nil)
		helper.AssertSuccess(t, w)

		w = helper.MakeRequest(t, "POST", "/api/v1/apps", map[string]interface{}{
			"name":        "测试应用4",
			"code":        "test_app4",
			"description": "这是测试应用4",
		})
		helper.AssertError(t, w, http.StatusForbidden)
	})
}
