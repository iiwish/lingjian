package test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskFlow(t *testing.T) {
	helper := NewTestHelper(t)

	var taskID uint

	// 1. 测试创建SQL任务
	t.Run("创建SQL任务", func(t *testing.T) {
		taskData := map[string]interface{}{
			"app_id": 1,
			"name":   "测试SQL任务",
			"type":   "sql",
			"cron":   "*/5 * * * *",
			"content": map[string]interface{}{
				"sql": "SELECT COUNT(*) FROM users",
			},
			"timeout":     30,
			"retry_times": 3,
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks", taskData)
		resp := helper.AssertSuccess(t, w)
		data := resp["data"].(map[string]interface{})
		taskID = uint(data["id"].(float64))
	})

	// 2. 测试创建HTTP任务
	t.Run("创建HTTP任务", func(t *testing.T) {
		taskData := map[string]interface{}{
			"app_id": 1,
			"name":   "测试HTTP任务",
			"type":   "http",
			"cron":   "0 0 * * *",
			"content": map[string]interface{}{
				"url":    "https://api.example.com/test",
				"method": "POST",
				"headers": map[string]interface{}{
					"Content-Type": "application/json",
				},
			},
			"timeout":     60,
			"retry_times": 3,
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks", taskData)
		helper.AssertSuccess(t, w)
	})

	// 3. 测试更新任务
	t.Run("更新任务", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name": "更新后的SQL任务",
			"cron": "0 */1 * * *",
			"content": map[string]interface{}{
				"sql": "SELECT COUNT(*) FROM roles",
			},
			"timeout":     45,
			"retry_times": 5,
		}
		w := helper.MakeRequest(t, "PUT", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10), updateData)
		helper.AssertSuccess(t, w)
	})

	// 4. 测试执行任务
	t.Run("执行任务", func(t *testing.T) {
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/execute", nil)
		helper.AssertSuccess(t, w)

		// 等待任务执行完成
		time.Sleep(time.Second * 2)

		// 检查任务日志
		w = helper.MakeRequest(t, "GET", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/logs", nil)
		resp := helper.AssertSuccess(t, w)
		logs := resp["data"].([]interface{})
		assert.NotEmpty(t, logs)
	})

	// 5. 测试禁用任务
	t.Run("禁用任务", func(t *testing.T) {
		w := helper.MakeRequest(t, "PUT", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/status", map[string]int{"status": 0})
		helper.AssertSuccess(t, w)

		// 尝试执行已禁用的任务
		w = helper.MakeRequest(t, "POST", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/execute", nil)
		helper.AssertError(t, w, http.StatusBadRequest)
	})

	// 6. 测试创建元素触发器
	t.Run("创建元素触发器", func(t *testing.T) {
		triggerData := map[string]interface{}{
			"app_id":        1,
			"element_type":  "form",
			"element_id":    1,
			"trigger_point": "before",
			"type":          "sql",
			"content": map[string]interface{}{
				"sql": "SELECT * FROM users WHERE id = :id",
			},
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/triggers", triggerData)
		helper.AssertSuccess(t, w)
	})

	// 7. 测试错误处理
	t.Run("测试错误处理", func(t *testing.T) {
		// 测试创建无效类型的任务
		invalidTask := map[string]interface{}{
			"app_id": 1,
			"name":   "无效任务",
			"type":   "invalid",
			"cron":   "* * * * *",
			"content": map[string]interface{}{
				"data": "test",
			},
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/tasks", invalidTask)
		helper.AssertError(t, w, http.StatusBadRequest)

		// 测试创建危险SQL任务
		dangerousTask := map[string]interface{}{
			"app_id": 1,
			"name":   "危险SQL任务",
			"type":   "sql",
			"cron":   "* * * * *",
			"content": map[string]interface{}{
				"sql": "DROP TABLE users",
			},
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/tasks", dangerousTask)
		helper.AssertError(t, w, http.StatusBadRequest)

		// 测试创建无效HTTP任务
		invalidHTTPTask := map[string]interface{}{
			"app_id": 1,
			"name":   "无效HTTP任务",
			"type":   "http",
			"cron":   "* * * * *",
			"content": map[string]interface{}{
				"url": "invalid-url",
			},
		}
		w = helper.MakeRequest(t, "POST", "/api/v1/tasks", invalidHTTPTask)
		helper.AssertError(t, w, http.StatusBadRequest)
	})

	// 8. 测试并发控制
	t.Run("测试并发控制", func(t *testing.T) {
		// 启用任务
		w := helper.MakeRequest(t, "PUT", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/status", map[string]int{"status": 1})
		helper.AssertSuccess(t, w)

		// 并发执行同一个任务
		done := make(chan bool)
		go func() {
			w := helper.MakeRequest(t, "POST", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/execute", nil)
			if w.Code == http.StatusOK {
				done <- true
			} else {
				done <- false
			}
		}()

		// 第二个请求应该失败
		w = helper.MakeRequest(t, "POST", "/api/v1/tasks/"+strconv.FormatUint(uint64(taskID), 10)+"/execute", nil)
		helper.AssertError(t, w, http.StatusBadRequest)

		// 等待第一个请求完成
		success := <-done
		assert.True(t, success)
	})
}
