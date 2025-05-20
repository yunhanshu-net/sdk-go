package widget

import "github.com/yunhanshu-net/pkg/x/tagx"

// newSliderWidget 创建滑块组件
func newSliderWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	slider := &SliderWidget{
		Widget: WidgetSlider,
		Type:   TypeNumber,
	}

	tag := info.Tags
	if tag["min"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		slider.Min = 0
	}

	if tag["max"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		slider.Max = 100
	}

	if tag["step"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		slider.Step = 1
	}

	if tag["show_stops"] != "" {
		if tag["show_stops"] == "true" {
			slider.ShowStops = true
		}
	}

	if tag["range"] != "" {
		if tag["range"] == "true" {
			slider.Range = true
			slider.Type = TypeArray
		}
	}

	if tag["show_input"] != "" {
		if tag["show_input"] == "true" {
			slider.ShowInput = true
		}
	}

	if tag["default_value"] != "" {
		slider.DefaultValue = tag["default_value"]
	}

	if tag["uint"] != "" {
		slider.Unit = tag["uint"]
	}

	return slider, nil
}

// SliderWidget 滑块组件
type SliderWidget struct {
	// 组件类型，固定为slider
	Widget string `json:"widget"`
	// 数据类型，一般为number或array
	Type string `json:"type"`
	// 最小值
	Min int `json:"min,omitempty"`
	// 最大值
	Max int `json:"max,omitempty"`
	// 步长
	Step int `json:"step,omitempty"`
	// 是否显示间断点
	ShowStops bool `json:"show_stops,omitempty"`
	// 是否为范围选择
	Range bool `json:"range,omitempty"`
	// 是否显示输入框
	ShowInput bool `json:"show_input,omitempty"`
	// 是否显示tooltip
	ShowTooltip bool `json:"show_tooltip,omitempty"`
	// 默认值
	DefaultValue interface{} `json:"default_value,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
	// 刻度标记
	Marks map[int]string `json:"marks,omitempty"`
	// 格式化tooltip
	FormatTooltip string `json:"format_tooltip,omitempty"`

	//单位，MB，人民币，等等
	Unit string `json:"unit,omitempty"`
}

func (w *SliderWidget) GetValueType() string {
	return w.Type
}

func (w *SliderWidget) GetWidgetType() string {
	return w.Widget
}
