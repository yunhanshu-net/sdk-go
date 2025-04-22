package response

type Widget interface {
}

type InputWidget struct {
	//输入框的类型，常规输入框还是文本域
	//LineText(默认),TextArea
	InputMode string `json:"input_mode,omitempty"`
	//文本限制输入框的文本限制ps:1-100
	InputTextLimit string `json:"input_text_limit,omitempty"`
	//占位符（文本框提示信息）
	InputPlaceholder string `json:"input_placeholder,omitempty"`
	//数值限制ps：1-100
	InputNumberLimit string `json:"input_number_limit,omitempty"`
	//默认值
	DefaultValue string `json:"default_value,omitempty"`
}

type FuncParam struct {
	//英文标识
	Code string `json:"code,omitempty"`
	//中文名称
	Name string `json:"name,omitempty"`
	//中文介绍
	Desc string `json:"desc,omitempty"`

	//数据类型，string，number，time，float，select
	Type string `json:"type,omitempty"`

	//会把该字段渲染成什么组件，默认：input
	//input(文本输入框,请求参数默认是输入框，响应参数默认是文本展示框）
	//checkbox(多选框)
	//radio(单选框)
	//select(下拉框)
	//switch(开关)
	//slider(滑块)
	//file(文件上传组件)
	Widget string `json:"widget,omitempty"`

	//是否必填
	Required string `json:"required,omitempty"`

	//示例数据
	Example string `json:"example,omitempty"`

	//多选框的选项，,逗号分割
	Options string `json:"options,omitempty"`

	InputWidget

	SelectOptions string `json:"select_options,omitempty"`
	FileSizeLimit string `json:"file_size_limit,omitempty"`
	FileTypeLimit string `json:"file_type_limit,omitempty"`
	IsTableField  bool   `json:"is_table_field"`
}

type ApiInfo struct {
	Router      string   `json:"router"`
	Method      string   `json:"method"`
	ApiDesc     string   `json:"api_desc"`
	Labels      []string `json:"labels"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	Classify    string   `json:"classify"`

	ParamsIn  []FuncParam `json:"params_in"`
	ParamsOut []FuncParam `json:"params_out"`
	UseTables []string    `json:"use_tables"`

	Callbacks []string `json:"callbacks"`
}
