package widget

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"github.com/yunhanshu-net/sdk-go/view/widget/types"
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
	//input
	//输入框的类型，常规输入框还是文本域
	//line_text(默认),text_area
	Mode string `json:"mode"`
	//占位符（文本框提示信息）
	Placeholder string `json:"placeholder"`
	//默认值
	DefaultValue string `json:"default_value"`
}

// NewInputWidget 创建输入框组件
func NewInputWidget(info *tagx.RunnerFieldInfo) (*InputWidget, error) {
	valueType := info.GetValueType()
	if !types.IsValueType(valueType) {
		return nil, fmt.Errorf("不是合法的值类型：%s", valueType)
	}
	if info == nil {
		return nil, fmt.Errorf("<UNK>nil")
	}
	if info.Tags == nil {
		info.Tags = map[string]string{}
	}

	input := &InputWidget{
		Mode:         stringsx.DefaultString(info.Tags["mode"], "line_text"),
		Placeholder:  info.Tags["placeholder"],
		DefaultValue: info.Tags["default_value"],
	}
	return input, nil
}

func (w *InputWidget) GetWidgetType() string {
	return types.WidgetInput
}
