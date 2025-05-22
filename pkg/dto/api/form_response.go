package api

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/view/widget"
	"reflect"
	"strings"
)

type FormResponseParamInfo struct {
	//英文标识
	Code string `json:"code"`
	//中文名称
	Name string `json:"name"`
	//中文介绍
	Desc string `json:"desc"`
	//是否必填
	Required bool `json:"required"`

	Callbacks    string      `json:"callbacks"`
	Validates    string      `json:"validates"`
	WidgetConfig interface{} `json:"widget_config"` //这里是widget.Widget类型的接口
	WidgetType   string      `json:"widget_type"`
	ValueType    string      `json:"value_type"`
	Example      string      `json:"example"`
}

type FormResponseParams struct {
	SearchCondList string `json:"search_cond_list"` //支持的查询条件

	RenderType string                   `json:"render_type"`
	Children   []*FormResponseParamInfo `json:"children"`
}

func (p *FormResponseParams) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

func newFormResponseParamInfo(tag *tagx.RunnerFieldInfo) (*FormResponseParamInfo, error) {

	widgetIns, err := widget.NewWidget(tag, response.RenderTypeForm)
	if err != nil {
		return nil, err
	}

	param := &FormResponseParamInfo{
		Code:         tag.GetCode(),
		Name:         tag.GetName(),
		Desc:         tag.GetDesc(),
		Required:     tag.GetRequired(),
		Validates:    tag.GetValidates(),
		Callbacks:    tag.GetCallbacks(),
		WidgetConfig: widgetIns,
		WidgetType:   widgetIns.GetWidgetType(),
		ValueType:    tag.GetValueType(),
		Example:      tag.GetExample(),
	}

	return param, nil
}

func NewFormResponseParams(el interface{}) (*FormResponseParams, error) {
	//renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	rspType := reflect.TypeOf(el)
	if rspType.Kind() == reflect.Pointer {
		rspType = rspType.Elem()
	}
	if rspType.Kind() != reflect.Struct && rspType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("输出参数仅支持Struct和Slice类型")
	}

	var tags []*tagx.RunnerFieldInfo
	if rspType.Kind() == reflect.Struct {
		resFields, err := tagx.ParseStructFieldsTypeOf(rspType, "runner")
		if err != nil {
			return nil, err
		}
		tags = resFields
	} else {
		tp, err := tagx.GetSliceElementType(el)
		if err != nil {
			return nil, err
		}
		of, err := tagx.ParseStructFieldsTypeOf(tp, "runner")
		if err != nil {
			return nil, err
		}
		tags = of
	}
	var searchCond []string
	children := make([]*FormResponseParamInfo, 0, len(tags))
	for _, field := range tags {
		if field.IsSearchCond() {
			searchCond = append(searchCond, field.GetCode())
			continue
		}
		info, err := newFormResponseParamInfo(field)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}
	return &FormResponseParams{SearchCondList: strings.Join(searchCond, ","), RenderType: response.RenderTypeForm, Children: children}, nil
}
