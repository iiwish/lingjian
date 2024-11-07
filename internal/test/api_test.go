package test

import (
	"testing"
)

func TestAPIFlow(t *testing.T) {
	helper := NewTestHelper(t)

	// 1. 测试用户登录
	t.Run("用户登录", func(t *testing.T) {
		loginData := map[string]string{
			"username": "admin",
			"password": "admin123",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/auth/login", loginData)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].(map[string]interface{})
		if token, ok := data["token"].(string); ok {
			helper.Token = token
		}
	})

	// 2. 测试刷新token
	t.Run("刷新token", func(t *testing.T) {
		w := helper.MakeRequest(t, "POST", "/api/v1/auth/refresh", nil)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].(map[string]interface{})
		if token, ok := data["token"].(string); ok {
			helper.Token = token
		}
	})

	// 3. 测试创建应用
	t.Run("创建应用", func(t *testing.T) {
		appData := map[string]interface{}{
			"name":        "测试应用",
			"code":        "test_app",
			"description": "这是一个测试应用",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/apps", appData)
		helper.AssertSuccess(t, w)
	})

	// 4. 测试创建角色
	t.Run("创建角色", func(t *testing.T) {
		roleData := map[string]interface{}{
			"name":     "应用管理员",
			"code":     "app_admin",
			"app_code": "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
		helper.AssertSuccess(t, w)
	})

	// 5. 测试分配权限
	t.Run("分配权限", func(t *testing.T) {
		permData := map[string]interface{}{
			"permission_codes": []string{"view_users", "create_user"},
			"app_code":         "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles/app_admin/permissions", permData)
		helper.AssertSuccess(t, w)
	})

	// 6. 测试创建定时任务
	t.Run("创建定时任务", func(t *testing.T) {
		taskData := map[string]interface{}{
			"name":        "数据清理任务",
			"cron":        "0 0 * * *",
			"sql":         "DELETE FROM logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY)",
			"description": "每天零点清理7天前的日志",
			"app_code":    "test_app",
			"triggers": []map[string]interface{}{
				{
					"type": "before",
					"sql":  "SET @start_time = NOW()",
				},
				{
					"type": "after",
					"sql":  "INSERT INTO task_logs (task_id, duration) VALUES (@task_id, TIMESTAMPDIFF(SECOND, @start_time, NOW()))",
				},
			},
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks", taskData)
		helper.AssertSuccess(t, w)
	})

	// 7. 测试获取应用列表
	t.Run("获取应用列表", func(t *testing.T) {
		w := helper.MakeRequest(t, "GET", "/api/v1/apps", nil)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].([]interface{})

		// 验证返回的应用列表中包含我们创建的应用
		found := false
		for _, app := range data {
			appMap := app.(map[string]interface{})
			if appMap["code"].(string) == "test_app" {
				found = true
				break
			}
		}
		if !found {
			t.Error("未找到创建的测试应用")
		}
	})
}
