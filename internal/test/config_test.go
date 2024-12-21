package test

import (
	"strconv"
	"testing"
)

func TestConfigTableItemFlow(t *testing.T) {
	helper := NewTestHelper(t)

	var tableID uint
	var itemID uint

	// 1. 创建测试表
	t.Run("创建数据表", func(t *testing.T) {
		tableData := map[string]interface{}{
			"name": "测试表",
			"code": "test_table",
			"fields": []map[string]interface{}{
				{
					"name":     "name",
					"type":     "varchar",
					"length":   50,
					"required": true,
				},
				{
					"name":     "age",
					"type":     "int",
					"required": false,
				},
			},
		}
		w := helper.MakeRequest(t, "POST", "/api/v1/config/tables", tableData)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				tableID = uint(id)
				t.Logf("创建的表ID: %d", tableID)
			}
		}
	})

	// 2. 创建表记录
	t.Run("创建表记录", func(t *testing.T) {
		itemData := map[string]interface{}{
			"name": "张三",
			"age":  25,
		}
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items"
		w := helper.MakeRequest(t, "POST", path, itemData)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				itemID = uint(id)
				t.Logf("创建的记录ID: %d", itemID)
			}
		}
	})

	// 3. 获取表记录
	t.Run("获取表记录", func(t *testing.T) {
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items/" + strconv.FormatUint(uint64(itemID), 10)
		w := helper.MakeRequest(t, "GET", path, nil)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			name, ok := data["name"].(string)
			if !ok || name != "张三" {
				t.Error("获取的记录数据不正确")
			}
			age, ok := data["age"].(float64)
			if !ok || age != 25 {
				t.Error("获取的记录数据不正确")
			}
		}
	})

	// 4. 更新表记录
	t.Run("更新表记录", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name": "李四",
			"age":  30,
		}
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items/" + strconv.FormatUint(uint64(itemID), 10)
		w := helper.MakeRequest(t, "PUT", path, updateData)
		helper.AssertSuccess(t, w)

		// 验证更新结果
		w = helper.MakeRequest(t, "GET", path, nil)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			name, ok := data["name"].(string)
			if !ok || name != "李四" {
				t.Error("更新后的记录数据不正确")
			}
			age, ok := data["age"].(float64)
			if !ok || age != 30 {
				t.Error("更新后的记录数据不正确")
			}
		}
	})

	// 5. 批量创建记录
	t.Run("批量创建记录", func(t *testing.T) {
		batchData := []map[string]interface{}{
			{
				"name": "王五",
				"age":  35,
			},
			{
				"name": "赵六",
				"age":  40,
			},
		}
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items/batch"
		w := helper.MakeRequest(t, "POST", path, batchData)
		helper.AssertSuccess(t, w)
	})

	// 6. 获取记录列表
	t.Run("获取记录列表", func(t *testing.T) {
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items"
		w := helper.MakeRequest(t, "GET", path, nil)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			items, ok := data["items"].([]interface{})
			if !ok {
				t.Error("获取的列表数据格式不正确")
			}
			if len(items) != 3 { // 1个单独创建 + 2个批量创建
				t.Errorf("获取的记录数量不正确,期望3条,实际%d条", len(items))
			}
		}
	})

	// 7. 带条件的列表查询

	// 8. 批量删除记录
	t.Run("批量删除记录", func(t *testing.T) {
		deleteIDs := []uint{itemID} // 删除第一条记录
		path := "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items/batch"
		w := helper.MakeRequest(t, "DELETE", path, deleteIDs)
		helper.AssertSuccess(t, w)

		// 验证删除结果
		path = "/api/v1/config/tables/" + strconv.FormatUint(uint64(tableID), 10) + "/items"
		w = helper.MakeRequest(t, "GET", path, nil)
		resp := helper.AssertSuccess(t, w)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			items, ok := data["items"].([]interface{})
			if !ok {
				t.Error("获取的列表数据格式不正确")
			}
			if len(items) != 2 { // 应该剩下2条记录
				t.Errorf("删除后的记录数量不正确,期望2条,实际%d条", len(items))
			}
		}
	})
}
