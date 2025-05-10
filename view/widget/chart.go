package widget

// ChartWidget 图表组件
type ChartWidget struct {
	// 组件类型，可选值：line(折线图)、bar(柱状图)、pie(饼图)、scatter(散点图)、radar(雷达图)
	Widget string `json:"widget"`
	// 数据类型，一般为object
	Type string `json:"type"`
	// 图表标题
	Title string `json:"title,omitempty"`
	// 图表子标题
	Subtitle string `json:"subtitle,omitempty"`
	// X轴数据
	XAxis []string `json:"x_axis,omitempty"`
	// Y轴数据
	YAxis []interface{} `json:"y_axis,omitempty"`
	// 系列名称
	SeriesNames []string `json:"series_names,omitempty"`
	// 图表高度
	Height string `json:"height,omitempty"`
	// 是否显示图例
	ShowLegend bool `json:"show_legend,omitempty"`
	// 是否显示工具箱
	ShowToolbox bool `json:"show_toolbox,omitempty"`
	// 是否支持缩放
	Zoomable bool `json:"zoomable,omitempty"`
	// 主题：light, dark
	Theme string `json:"theme,omitempty"`
}

func (w *ChartWidget) GetValueType() string {
	return w.Type
}

func (w *ChartWidget) GetWidgetType() string {
	return w.Widget
}
