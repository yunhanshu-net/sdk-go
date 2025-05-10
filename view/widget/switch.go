package widget

import "github.com/yunhanshu-net/sdk-go/pkg/tagx"

// SwitchWidget 开关组件
type SwitchWidget struct {
	// 组件类型，固定为switch
	Widget string `json:"widget"`
	// 数据类型，一般为boolean
	Type string `json:"type"`
	// 打开时的文字
	ActiveText string `json:"active_text,omitempty"`
	// 关闭时的文字
	InactiveText string `json:"inactive_text,omitempty"`
	// 打开时的值
	ActiveValue interface{} `json:"active_value,omitempty"`
	// 关闭时的值
	InactiveValue interface{} `json:"inactive_value,omitempty"`
	// 打开时的颜色
	ActiveColor string `json:"active_color,omitempty"`
	// 关闭时的颜色
	InactiveColor string `json:"inactive_color,omitempty"`
	// 默认值
	DefaultValue interface{} `json:"default_value,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
	// 尺寸：large, default, small
	Size string `json:"size,omitempty"`
	// 是否显示内部文字
	InlinePrompt bool `json:"inline_prompt,omitempty"`
}

// newSwitchWidget 创建开关组件
func newSwitchWidget(info *tagx.FieldInfo) (Widget, error) {
	switch_ := &SwitchWidget{
		Widget: WidgetSwitch,
		Type:   TypeBoolean,
	}

	tag := info.Tags
	if tag["active_text"] != "" {
		switch_.ActiveText = tag["active_text"]
	}

	if tag["inactive_text"] != "" {
		switch_.InactiveText = tag["inactive_text"]
	}

	if tag["active_color"] != "" {
		switch_.ActiveColor = tag["active_color"]
	}

	if tag["inactive_color"] != "" {
		switch_.InactiveColor = tag["inactive_color"]
	}

	if tag["default_value"] != "" {
		switch_.DefaultValue = tag["default_value"] == "true"
	}

	if tag["size"] != "" {
		switch_.Size = tag["size"]
	}

	return switch_, nil
}

func (w *SwitchWidget) GetValueType() string {
	return w.Type
}

func (w *SwitchWidget) GetWidgetType() string {
	return w.Widget
}
