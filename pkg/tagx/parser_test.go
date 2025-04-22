package tagx

import (
	"reflect"
	"testing"
)

// 测试数据结构定义
type BasicForm struct {
	Username string `form:"code:username;name:用户名;desc:登录账户;required:true;example:zhangsan;default:guest;type:string;widget:input"`
}

type TextBoxForm struct {
	Comment string `form:"placeholder:请输入评论...;fuzzy:true;max:500;min:10"`
	User    string `form:"required"`
}

type NestedForm struct {
	BasicForm `form:"code:base"` // 匿名字段
	Email     string             `form:"code:email;type:email;widget:input"`
}

// 测试用例
func TestParseStructFields(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		tagKey   string
		expected []FieldInfo
	}{
		{
			name:   "基础字段解析",
			input:  BasicForm{},
			tagKey: "form",
			expected: []FieldInfo{
				{
					Name: "Username",
					Tags: map[string]string{
						"code":     "username",
						"name":     "用户名",
						"desc":     "登录账户",
						"required": "true",
						"example":  "zhangsan",
						"default":  "guest",
						"type":     "string",
						"widget":   "input",
					},
					Type: reflect.TypeOf(""),
				},
			},
		},
		{
			name:   "文本框属性解析",
			input:  TextBoxForm{},
			tagKey: "form",
			expected: []FieldInfo{
				{
					Name: "Comment",
					Tags: map[string]string{
						"placeholder": "请输入评论...",
						"fuzzy":       "true",
						"max":         "500",
						"min":         "10",
					},
					Type: reflect.TypeOf(""),
				},
				{
					Name: "User",
					Tags: map[string]string{
						"required": "",
					},
					Type: reflect.TypeOf(""),
				},
			},
		},
		{
			name:   "嵌套匿名字段解析",
			input:  NestedForm{},
			tagKey: "form",
			expected: []FieldInfo{
				{
					Name: "BasicForm.Username", // 匿名字段层级
					Tags: map[string]string{
						"code": "username",
						"name": "用户名",
						// ...其他标签同上
					},
					Type: reflect.TypeOf(""),
				},
				{
					Name: "Email",
					Tags: map[string]string{
						"code":   "email",
						"type":   "email",
						"widget": "input",
					},
					Type: reflect.TypeOf(""),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseStructFields(tt.input, tt.tagKey)
			if err != nil {
				t.Fatalf("解析失败: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("字段数量不符: 期望 %d, 实际 %d",
					len(tt.expected), len(result))
			}

			for i, field := range result {
				exp := tt.expected[i]

				// 验证字段名和类型
				//if field.Name != exp.Name || field.Type != exp.Type {
				//	t.Errorf("字段 %d 不匹配:\n期望 %s (%s)\n实际 %s (%s)",
				//		i, exp.Name, exp.Type, field.Name, field.Type)
				//}

				// 验证标签键值对
				for k, v := range exp.Tags {
					actualVal, ok := field.Tags[k]
					if !ok || actualVal != v {
						t.Errorf("标签 %s 错误:\n期望 %s\n实际 %s", k, v, actualVal)
					}
				}
			}
		})
	}
}

// 错误场景测试
func TestParseErrors(t *testing.T) {
	t.Run("非结构体输入", func(t *testing.T) {
		_, err := ParseStructFields("invalid", "form")
		if err == nil {
			t.Error("预期报错但未发生")
		}
	})

	t.Run("无效标签格式", func(t *testing.T) {
		type InvalidForm struct {
			Test string `form:"code=123"` // 错误的分隔符
		}
		_, err := ParseStructFields(InvalidForm{}, "form")
		if err == nil {
			t.Error("预期标签解析错误但未发生")
		}
	})
}
