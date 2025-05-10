package widget

import (
	"errors"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"github.com/yunhanshu-net/sdk-go/view/widget/types"
	"strings"
)

// SelectWidget 下拉框组件
type SelectWidget struct {
	// 选项列表
	Options []string `json:"options"`
	// 是否多选
	Multiple bool `json:"multiple,omitempty"`

	// 默认值
	DefaultValue string `json:"default_value,omitempty"`
}

// NewSelectWidget 创建下拉框组件
func NewSelectWidget(info *tagx.FieldInfo) (Widget, error) {
	if info == nil {
		return nil, errors.New("<UNK>")
	}
	if info.Tags == nil {
		info.Tags = make(map[string]string)
	}
	select_ := &SelectWidget{
		DefaultValue: info.Tags["default_value"],
		Options:      strings.Split(info.Tags["options"], ","),
	}
	tag := info.Tags
	if tag["options"] != "" {
		select_.Options = strings.Split(tag["options"], ",")
	}
	if _, ok := tag["multiple"]; ok {
		select_.Multiple = true
	}

	return select_, nil
}

func (w *SelectWidget) GetWidgetType() string {
	return types.WidgetSelect
}
