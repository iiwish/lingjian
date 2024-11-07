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
			name:        "查看用户列表",
			code:        "view_users",
			appCode:     "test_app1",
			type_:       "api",
			path:        "/api/v1/users",
			method:      "GET",
			description: "允许查看用户列表",
		},
		{
			name:        "创建用户",
			code:        "create_user",
			appCode:     "test_app1",
			type_:       "api",
			path:        "/api/v1/users",
			method:      "POST",
			description: "允许创建新用户",
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
		permCodes := []string{"view_users", "create_user"}
		assignData := map[string]interface{}{
			"permission_codes": permCodes,
			"app_code":         "test_app1",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles/super_admin/permissions", assignData)
		helper.AssertSuccess(t, w)

		// 普通管理员只有查看权限
		assignData = map[string]interface{}{
			"permission_codes": []string{"view_users"},
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
		permissions := resp["data"].([]interface{})
		assert.Equal(t, 2, len(permissions))

		// 检查普通管理员权限
		w = helper.MakeRequest(t, "GET", "/api/v1/roles/admin/permissions?app_code=test_app1", nil)
		resp = helper.AssertSuccess(t, w)
		permissions = resp["data"].([]interface{})
		assert.Equal(t, 1, len(permissions))
	})

	// 6. 测试权限检查
	t.Run("测试权限检查", func(t *testing.T) {
		// 超级管理员可以访问所有接口
		w := helper.MakeRequest(t, "GET", "/api/v1/users", nil)
		helper.AssertSuccess(t, w)

		w = helper.MakeRequest(t, "POST", "/api/v1/users", map[string]string{"username": "test"})
		helper.AssertSuccess(t, w)

		// 切换到普通管理员角色
		w = helper.MakeRequest(t, "POST", "/api/v1/auth/switch-role", map[string]string{"role_code": "admin"})
		helper.AssertSuccess(t, w)

		// 普通管理员只能查看
		w = helper.MakeRequest(t, "GET", "/api/v1/users", nil)
		helper.AssertSuccess(t, w)

		w = helper.MakeRequest(t, "POST", "/api/v1/users", map[string]string{"username": "test"})
		helper.AssertError(t, w, http.StatusForbidden)
	})

	// 7. 测试动态权限更新
	t.Run("测试动态权限更新", func(t *testing.T) {
		// 创建新权限
		newPermData := map[string]interface{}{
			"name":        "删除用户",
			"code":        "delete_user",
			"app_code":    "test_app1",
			"type":        "api",
			"path":        "/api/v1/users/{id}",
			"method":      "DELETE",
			"description": "允许删除用户",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/permissions", newPermData)
		helper.AssertSuccess(t, w)

		// 分配给超级管理员
		assignData := map[string]interface{}{
			"permission_codes": []string{"delete_user"},
			"app_code":         "test_app1",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/roles/super_admin/permissions", assignData)
		helper.AssertSuccess(t, w)

		// 验证权限是否生效
		w = helper.MakeRequest(t, "DELETE", "/api/v1/users/1", nil)
		helper.AssertSuccess(t, w)
	})

	// 8. 测试多应用隔离
	t.Run("测试多应用隔离", func(t *testing.T) {
		// 创建应用2的角色
		roleData := map[string]interface{}{
			"name":     "应用2管理员",
			"code":     "app2_admin",
			"app_code": "test_app2",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
		helper.AssertSuccess(t, w)

		// 创建应用2的权限
		permData := map[string]interface{}{
			"name":        "应用2查看权限",
			"code":        "app2_view",
			"app_code":    "test_app2",
			"type":        "api",
			"path":        "/api/v1/app2/data",
			"method":      "GET",
			"description": "应用2的查看权限",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/permissions", permData)
		helper.AssertSuccess(t, w)

		// 分配权限
		assignData := map[string]interface{}{
			"permission_codes": []string{"app2_view"},
			"app_code":         "test_app2",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/roles/app2_admin/permissions", assignData)
		helper.AssertSuccess(t, w)

		// 验证应用1的角色无法访问应用2的资源
		w = helper.MakeRequest(t, "GET", "/api/v1/app2/data", nil)
		helper.AssertError(t, w, http.StatusForbidden)
	})
}

func TestRBACErrorCases(t *testing.T) {
	helper := NewTestHelper(t)

	// 1. 测试创建重复角色
	t.Run("创建重复角色", func(t *testing.T) {
		roleData := map[string]interface{}{
			"name":     "测试角色",
			"code":     "test_role",
			"app_code": "test_app",
		}

		// 第一次创建
		w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
		helper.AssertSuccess(t, w)

		// 第二次创建同名角色
		w = helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
		helper.AssertError(t, w, http.StatusConflict)
	})

	// 2. 测试创建无效权限
	t.Run("创建无效权限", func(t *testing.T) {
		permData := map[string]interface{}{
			"name":        "无效权限",
			"code":        "invalid_perm",
			"app_code":    "test_app",
			"type":        "invalid_type", // 无效的权限类型
			"path":        "/api/v1/test",
			"method":      "GET",
			"description": "测试无效权限",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/permissions", permData)
		helper.AssertError(t, w, http.StatusBadRequest)
	})

	// 3. 测试分配不存在的权限
	t.Run("分配不存在的权限", func(t *testing.T) {
		assignData := map[string]interface{}{
			"permission_codes": []string{"non_existent_perm"},
			"app_code":         "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles/test_role/permissions", assignData)
		helper.AssertError(t, w, http.StatusNotFound)
	})

	// 4. 测试循环角色继承
	t.Run("测试循环角色继承", func(t *testing.T) {
		// 创建角色A
		roleAData := map[string]interface{}{
			"name":     "角色A",
			"code":     "role_a",
			"app_code": "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleAData)
		helper.AssertSuccess(t, w)

		// 创建角色B，继承自角色A
		roleBData := map[string]interface{}{
			"name":        "角色B",
			"code":        "role_b",
			"app_code":    "test_app",
			"parent_code": "role_a",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/roles", roleBData)
		helper.AssertSuccess(t, w)

		// 尝试让角色A继承角色B，这会造成循环
		updateData := map[string]interface{}{
			"parent_code": "role_b",
		}
		w = helper.MakeRequest(t, "PUT", "/api/v1/roles/role_a", updateData)
		helper.AssertError(t, w, http.StatusBadRequest)
	})

	// 5. 测试跨应用权限分配
	t.Run("测试跨应用权限分配", func(t *testing.T) {
		// 创建应用1的权限
		permData := map[string]interface{}{
			"name":        "应用1权限",
			"code":        "app1_perm",
			"app_code":    "app1",
			"type":        "api",
			"path":        "/api/v1/test",
			"method":      "GET",
			"description": "应用1的测试权限",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/permissions", permData)
		helper.AssertSuccess(t, w)

		// 尝试将应用1的权限分配给应用2的角色
		assignData := map[string]interface{}{
			"permission_codes": []string{"app1_perm"},
			"app_code":         "app2",
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/roles/app2_role/permissions", assignData)
		helper.AssertError(t, w, http.StatusBadRequest)
	})
}
