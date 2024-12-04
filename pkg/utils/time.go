package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

// MarshalJSON 实现自定义的 JSON 序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format(ctLayout))
}

// UnmarshalJSON 实现自定义的 JSON 反序列化
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]
	t, err := time.Parse(ctLayout, str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Scan 实现 sql.Scanner 接口
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime{Time: v}
		return nil
	case []byte:
		t, err := time.Parse(ctLayout, string(v))
		if err != nil {
			return err
		}
		*ct = CustomTime{Time: t}
		return nil
	case string:
		t, err := time.Parse(ctLayout, v)
		if err != nil {
			return err
		}
		*ct = CustomTime{Time: t}
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *utils.CustomTime", value)
	}
}

// Value 实现 driver.Valuer 接口
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// NewCustomTime 创建一个 CustomTime 实例
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

// NowCustomTime 返回当前时间的 CustomTime 格式
func NowCustomTime() CustomTime {
	return NewCustomTime(time.Now())
}
