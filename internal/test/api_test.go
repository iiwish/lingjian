package test

import (
	"testing"
)

func TestAPIFlow(t *testing.T) {
	helper := NewTestHelper(t)

	var appID uint

	// 1. 测试用户登录 - 已经在NewTestHelper中完成

	// 2. 测试刷新token
	t.Run("刷新token", func(t *testing.T) {
		headers := map[string]string{
			"X-Refresh-Token": helper.RefreshToken,
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/auth/refresh", nil, headers)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].(map[string]interface{})
		if token, ok := data["access_token"].(string); ok {
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
		// 添加错误响应日志
		if w.Code != 200 {
			t.Logf("创建应用失败，响应状态码: %d", w.Code)
			t.Logf("响应内容: %s", w.Body.String())
		}
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				appID = uint(id)
				t.Logf("创建的应用ID: %d", appID)
			}
		}
	})

	// 4. 测试创建角色
	t.Run("创建角色", func(t *testing.T) {
		roleData := map[string]interface{}{
			"name":     "应用管理员",
			"code":     "app_admin",
			"app_code": "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles", roleData)
		// 添加错误响应日志
		if w.Code != 200 {
			t.Logf("创建角色失败，响应状态码: %d", w.Code)
			t.Logf("响应内容: %s", w.Body.String())
		}
		helper.AssertSuccess(t, w)
	})

	// 5. 测试分配权限
	t.Run("分配权限", func(t *testing.T) {
		permData := map[string]interface{}{
			"permission_codes": []string{"view_apps", "create_app", "create_role", "assign_permission", "create_task"},
			"app_code":         "test_app",
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/roles/app_admin/permissions", permData)
		// 添加错误响应日志
		if w.Code != 200 {
			t.Logf("分配权限失败，响应状态码: %d", w.Code)
			t.Logf("响应内容: %s", w.Body.String())
		}
		helper.AssertSuccess(t, w)
	})

	// 6. 测试创建定时任务
	t.Run("创建定时任务", func(t *testing.T) {
		taskData := map[string]interface{}{
			"app_id": appID,
			"name":   "数据清理任务",
			"type":   "sql",
			"cron":   "0 0 * * *",
			"content": map[string]interface{}{
				"sql": "DELETE FROM sys_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY)",
			},
			"timeout":     60,
			"retry_times": 3,
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks/scheduled", taskData)
		// 添加错误响应日志
		if w.Code != 200 {
			t.Logf("创建定时任务失败，响应状态码: %d", w.Code)
			t.Logf("响应内容: %s", w.Body.String())
		}
		helper.AssertSuccess(t, w)
	})

	// 7. 测试获取应用列表
	t.Run("获取应用列表", func(t *testing.T) {
		w := helper.MakeRequest(t, "GET", "/api/v1/apps", nil)
		// 添加错误响应日志
		if w.Code != 200 {
			t.Logf("获取应用列表失败，响应状态码: %d", w.Code)
			t.Logf("响应内容: %s", w.Body.String())
		}
		resp := helper.AssertSuccess(t, w)

		// 安全地获取data字段
		data, ok := resp["data"]
		if !ok {
			t.Fatal("响应中没有data字段")
		}

		// 获取items字段
		items, ok := data.(map[string]interface{})["items"]
		if !ok {
			t.Fatal("data字段中没有items字段")
		}

		// 将items转换为[]interface{}
		appList, ok := items.([]interface{})
		if !ok {
			t.Fatal("items字段不是数组类型")
		}

		// 验证返回的应用列表中包含我们创建的应用
		found := false
		for _, app := range appList {
			appMap, ok := app.(map[string]interface{})
			if !ok {
				continue
			}
			code, ok := appMap["code"].(string)
			if !ok {
				continue
			}
			if code == "test_app" {
				found = true
				break
			}
		}
		if !found {
			t.Error("未找到创建的测试应用")
		}
	})
}
