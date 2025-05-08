package render

import (
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
)

// InputWidget 输入框组件
type InputWidget struct {
	//会把该字段渲染成什么组件，默认：input
	//input(文本输入框,请求参数默认是输入框，响应参数默认是文本展示框）
	//checkbox(多选框)
	//radio(单选框)
	//select(下拉框)
	//switch(开关)
	//slider(滑块)
	//file(文件上传组件)

	//input
	Widget string `json:"widget"`
	//数据类型，string，number，time，float
	Type string `json:"type"`
	//输入框的类型，常规输入框还是文本域
	//line_text(默认),text_area
	Mode string `json:"mode"`
	//文本限制输入框的文本限制ps:1-100
	TextLimit string `json:"text_limit"`
	//占位符（文本框提示信息）
	Placeholder string `json:"placeholder"`
	//数值限制ps：1-100
	NumberLimit string `json:"number_limit"`
	//默认值
	DefaultValue string `json:"default_value"`
	//示例数据
	Example string `json:"example"`
}

// newInputWidget 创建输入框组件
func newInputWidget(info *tagx.FieldInfo) (Widget, error) {
	input := &InputWidget{
		Widget:    WidgetInput,
		Type:      TypeString,
		Mode:      "line_text",
		TextLimit: info.Tags["text_limit"],
	}

	tag := info.Tags
	if tag["mode"] != "" {
		input.Mode = tag["mode"]
	}

	if tag["text_limit"] != "" {
		input.TextLimit = tag["text_limit"]
	}

	if tag["placeholder"] != "" {
		input.Placeholder = tag["placeholder"]
	}

	if tag["number_limit"] != "" {
		input.NumberLimit = tag["number_limit"]
	}
	if tag["default_value"] != "" {
		input.DefaultValue = tag["default_value"]
	}

	if tag["example"] != "" {
		input.Example = tag["example"]
	}
	return input, nil
}

func (w *InputWidget) GetType() string {
	return w.Type
}

func (w *InputWidget) GetWidget() string {
	return w.Widget
}
