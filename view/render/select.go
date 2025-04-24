package render

import "github.com/yunhanshu-net/sdk-go/pkg/tagx"

// SelectWidget 下拉框组件
type SelectWidget struct {
	// 组件类型，固定为select
	Widget string `json:"widget"`
	// 数据类型，一般为string或array
	Type string `json:"type"`
	// 选项列表
	Options []string `json:"options"`
	// 是否多选
	Multiple bool `json:"multiple,omitempty"`
	// 是否可搜索
	Filterable bool `json:"filterable,omitempty"`
	// 是否可清空
	Clearable bool `json:"clearable,omitempty"`
	// 占位符
	Placeholder string `json:"placeholder,omitempty"`
	// 默认值
	DefaultValue interface{} `json:"default_value,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
	// 尺寸：large, default, small
	Size string `json:"size,omitempty"`
	// 是否可创建新条目
	AllowCreate bool `json:"allow_create,omitempty"`
	// 是否显示全选选项
	ShowAllOption bool `json:"show_all_option,omitempty"`
	// 全选选项文本
	AllOptionLabel string `json:"all_option_label,omitempty"`
}

// newSelectWidget 创建下拉框组件
func newSelectWidget(info *tagx.FieldInfo) (Widget, error) {
	select_ := &SelectWidget{
		Widget: WidgetSelect,
		Type:   TypeString,
	}

	tag := info.Tags
	if tag["options"] != "" {
		select_.Options = []string{tag["options"]}
	}

	if tag["multiple"] != "" {
		if tag["multiple"] == "true" {
			select_.Multiple = true
			select_.Type = TypeArray
		}
	}

	if tag["filterable"] != "" {
		if tag["filterable"] == "true" {
			select_.Filterable = true
		}
	}

	if tag["clearable"] != "" {
		if tag["clearable"] == "true" {
			select_.Clearable = true
		}
	}

	if tag["placeholder"] != "" {
		select_.Placeholder = tag["placeholder"]
	}

	if tag["default_value"] != "" {
		select_.DefaultValue = tag["default_value"]
	}

	if tag["size"] != "" {
		select_.Size = tag["size"]
	}

	return select_, nil
}

func (w *SelectWidget) GetType() string {
	return w.Type
}

func (w *SelectWidget) GetWidget() string {
	return w.Widget
}
