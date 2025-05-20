package api

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/view/widget"
	"reflect"
)

type ResponseParamInfo struct {
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

type ResponseParams struct {
	RenderType string               `json:"render_type"`
	Children   []*ResponseParamInfo `json:"children"`
}

func (p *ResponseParams) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

func newResponseParamInfo(tag *tagx.RunnerFieldInfo, renderType string) (*ResponseParamInfo, error) {

	widgetIns, err := widget.NewWidget(tag, renderType)
	if err != nil {
		return nil, err
	}

	param := &ResponseParamInfo{
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

func NewResponseParams(el interface{}, renderType string) (*ResponseParams, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
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
	children := make([]*ResponseParamInfo, 0, len(tags))
	for _, field := range tags {
		info, err := newResponseParamInfo(field, renderType)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}
	return &ResponseParams{RenderType: renderType, Children: children}, nil
}
