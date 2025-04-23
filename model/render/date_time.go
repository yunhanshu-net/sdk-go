package render

// DateTimeWidget 日期时间选择器组件
type DateTimeWidget struct {
	// 组件类型，可选值：date(日期)、time(时间)、datetime(日期时间)、daterange(日期范围)、timerange(时间范围)
	Widget string `json:"widget"`
	// 数据类型，一般为time
	Type string `json:"type"`
	// 日期格式，如：yyyy-MM-dd、HH:mm:ss、yyyy-MM-dd HH:mm:ss
	Format string `json:"format,omitempty"`
	// 是否显示清除按钮
	Clearable bool `json:"clearable,omitempty"`
	// 占位符
	Placeholder string `json:"placeholder,omitempty"`
	// 默认值
	DefaultValue string `json:"default_value,omitempty"`
	// 最小日期/时间
	Min string `json:"min,omitempty"`
	// 最大日期/时间
	Max string `json:"max,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
}

func (w *DateTimeWidget) GetType() string {
	return w.Type
}

func (w *DateTimeWidget) GetWidget() string {
	return w.Widget
}
