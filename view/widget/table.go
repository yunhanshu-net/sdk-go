package widget

import (
	"github.com/pkg/errors"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

// TableWidget 表格组件
type TableWidget struct {
	// 组件类型，固定为table
	Widget string `json:"widget"`
	// 表格列定义
	Columns []TableColumn `json:"columns,omitempty"`
}

// TableColumn 表格列定义
type TableColumn struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	ValueType string `json:"value_type"`

	//// 列标识
	//Prop string `json:"prop"`
	//// 列标题
	//Label string `json:"label"`
	//// 列宽度
	//Width string `json:"width,omitempty"`
	//// 对齐方式：left, center, right
	//Align string `json:"align,omitempty"`
	//// 是否固定列：left, right
	//Fixed string `json:"fixed,omitempty"`
	//// 是否可排序
	//Sortable bool `json:"sortable,omitempty"`
	//// 格式化函数名称
	//Formatter string `json:"formatter,omitempty"`
}

func (w *TableWidget) GetValueType() string {
	return TypeArray
}

func (w *TableWidget) GetWidgetType() string {
	return w.Widget
}

func NewTable(info []*tagx.RunnerFieldInfo) (*TableWidget, error) {
	Columns := make([]TableColumn, 0)
	for _, v := range info {
		Columns = append(Columns, TableColumn{
			Code:      v.GetCode(),
			Name:      v.GetName(),
			ValueType: v.GetValueType(),
		})
	}
	if info == nil {
		return nil, errors.New("NewTable info ==nil")
	}

	return &TableWidget{Widget: WidgetTable, Columns: Columns}, nil

}
