package render

import "github.com/yunhanshu-net/sdk-go/pkg/tagx"

// CheckboxWidget 多选框组件
type CheckboxWidget struct {
	// 组件类型，固定为checkbox
	Widget string `json:"widget"`
	// 数据类型，一般为array
	Type string `json:"type"`
	// 选项列表
	Options []string `json:"options"`
	// 默认选中的值
	DefaultValue []string `json:"default_value,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
	// 按钮样式
	ButtonStyle bool `json:"button_style,omitempty"`
	// 尺寸：large, default, small
	Size string `json:"size,omitempty"`
	// 最小选中数量
	Min int `json:"min,omitempty"`
	// 最大选中数量
	Max int `json:"max,omitempty"`
}

// newCheckboxWidget 创建多选框组件
func newCheckboxWidget(info *tagx.FieldInfo) (Widget, error) {
	checkbox := &CheckboxWidget{
		Widget: WidgetCheckbox,
		Type:   TypeArray,
	}

	tag := info.Tags
	if tag["options"] != "" {
		checkbox.Options = []string{tag["options"]}
	}

	if tag["default_value"] != "" {
		checkbox.DefaultValue = []string{tag["default_value"]}
	}

	if tag["button_style"] != "" {
		if tag["button_style"] == "true" {
			checkbox.ButtonStyle = true
		}
	}

	if tag["size"] != "" {
		checkbox.Size = tag["size"]
	}

	if tag["min"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		checkbox.Min = 0
	}

	if tag["max"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		checkbox.Max = 0
	}

	return checkbox, nil
}

func (w *CheckboxWidget) GetType() string {
	return w.Type
}

func (w *CheckboxWidget) GetWidget() string {
	return w.Widget
}
