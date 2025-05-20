package widget

import "github.com/yunhanshu-net/pkg/x/tagx"

// RadioWidget 单选框组件
type RadioWidget struct {
	// 组件类型，固定为radio
	Widget string `json:"widget"`
	// 数据类型，一般为string
	Type string `json:"type"`
	// 选项列表
	Options []string `json:"options"`
	// 默认选中的值
	DefaultValue string `json:"default_value,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
	// 按钮样式
	ButtonStyle bool `json:"button_style,omitempty"`
	// 尺寸：large, default, small
	Size string `json:"size,omitempty"`
}

// newRadioWidget 创建单选框组件
func newRadioWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	radio := &RadioWidget{
		Widget: WidgetRadio,
		Type:   TypeString,
	}

	tag := info.Tags
	if tag["options"] != "" {
		radio.Options = []string{tag["options"]}
	}

	if tag["default_value"] != "" {
		radio.DefaultValue = tag["default_value"]
	}

	if tag["button_style"] != "" {
		if tag["button_style"] == "true" {
			radio.ButtonStyle = true
		}
	}

	if tag["size"] != "" {
		radio.Size = tag["size"]
	}

	return radio, nil
}

func (w *RadioWidget) GetValueType() string {
	return w.Type
}

func (w *RadioWidget) GetWidgetType() string {
	return w.Widget
}
