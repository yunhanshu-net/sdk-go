package render

type Params struct {

	//form,table,echarts,bi,3D .....
	Type     string       `json:"type"`
	Children []*ParamInfo `json:"children"`
}

type ParamInfo struct {
	//英文标识
	Code string `json:"code,omitempty"`
	//中文名称
	Name string `json:"name,omitempty"`
	//中文介绍
	Desc string `json:"desc,omitempty"`

	//是否必填
	Required string `json:"required,omitempty"`

	////多选框的选项，,逗号分割
	//Options string `json:"options,omitempty"`

	Widget Widget `json:"widget"`

	//SelectOptions string `json:"select_options,omitempty"`
	//FileSizeLimit string `json:"file_size_limit,omitempty"`
	//FileTypeLimit string `json:"file_type_limit,omitempty"`
	//IsTableField  bool   `json:"is_table_field"`
}

type ApiInfo struct {
	Router      string   `json:"router"`
	Method      string   `json:"method"`
	ApiDesc     string   `json:"api_desc"`
	Labels      []string `json:"labels"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	Classify    string   `json:"classify"`

	//输入参数
	ParamsIn Params `json:"params_in"`

	//输出参数
	ParamsOut Params   `json:"params_out"`
	UseTables []string `json:"use_tables"`

	Callbacks []string `json:"callbacks"`
}
