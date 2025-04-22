package render

type InputWidget struct {
	//会把该字段渲染成什么组件，默认：input
	//input(文本输入框,请求参数默认是输入框，响应参数默认是文本展示框）
	//checkbox(多选框)
	//radio(单选框)
	//select(下拉框)
	//switch(开关)
	//slider(滑块)
	//file(文件上传组件)
	Widget string `json:"widget,omitempty"`
	//数据类型，string，number，time，float
	Type string `json:"type"`
	//输入框的类型，常规输入框还是文本域
	//LineText(默认),TextArea
	Mode string `json:"mode,omitempty"`
	//文本限制输入框的文本限制ps:1-100
	TextLimit string `json:"text_limit,omitempty"`
	//占位符（文本框提示信息）
	Placeholder string `json:"placeholder,omitempty"`
	//数值限制ps：1-100
	NumberLimit string `json:"number_limit,omitempty"`
	//默认值
	DefaultValue string `json:"default_value,omitempty"`
	//示例数据
	Example string `json:"example,omitempty"`
}

func (w *InputWidget) GetType() string {
	return w.Type
}
func (w *InputWidget) GetWidget() string {
	return w.Widget
}
