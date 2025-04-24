package render

// TableWidget 表格组件
type TableWidget struct {
	// 组件类型，固定为table
	Widget string `json:"widget"`
	// 数据类型，一般为array
	Type string `json:"type"`
	// 表格列定义
	Columns []TableColumn `json:"columns,omitempty"`
	// 是否显示序号列
	ShowIndex bool `json:"show_index,omitempty"`
	// 是否显示选择框
	Selectable bool `json:"selectable,omitempty"`
	// 是否显示边框
	Border bool `json:"border,omitempty"`
	// 是否显示斑马纹
	Stripe bool `json:"stripe,omitempty"`
	// 是否支持排序
	Sortable bool `json:"sortable,omitempty"`
	// 是否支持过滤
	Filterable bool `json:"filterable,omitempty"`
	// 是否支持分页
	Pagination bool `json:"pagination,omitempty"`
	// 每页显示条数选项
	PageSizes []int `json:"page_sizes,omitempty"`
	// 默认每页显示条数
	DefaultPageSize int `json:"default_page_size,omitempty"`
}

// TableColumn 表格列定义
type TableColumn struct {
	// 列标识
	Prop string `json:"prop"`
	// 列标题
	Label string `json:"label"`
	// 列宽度
	Width string `json:"width,omitempty"`
	// 对齐方式：left, center, right
	Align string `json:"align,omitempty"`
	// 是否固定列：left, right
	Fixed string `json:"fixed,omitempty"`
	// 是否可排序
	Sortable bool `json:"sortable,omitempty"`
	// 格式化函数名称
	Formatter string `json:"formatter,omitempty"`
}

func (w *TableWidget) GetType() string {
	return w.Type
}

func (w *TableWidget) GetWidget() string {
	return w.Widget
}
