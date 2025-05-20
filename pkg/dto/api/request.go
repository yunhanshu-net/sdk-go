package api

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/view/widget"
	"reflect"
	"strings"
)

type RequestParamInfo struct {
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

type RequestParams struct {
	SearchCondList string `json:"search_cond_list"` //支持的查询条件
	//SearchCondBlickList map[string]string   `json:"search_cond_blick_list"` //禁止的查询条件
	RenderType string              `json:"render_type"`
	Children   []*RequestParamInfo `json:"children"`
}

func (p *RequestParams) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

func newRequestParamInfo(tag *tagx.RunnerFieldInfo, renderType string) (*RequestParamInfo, error) {

	widgetIns, err := widget.NewWidget(tag, renderType)
	if err != nil {
		return nil, err
	}
	param := &RequestParamInfo{
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

func NewRequestParams(el interface{}, renderType string) (*RequestParams, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	typeOf := reflect.TypeOf(el)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("输入参数仅支持Struct类型")
	}
	reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
	if err != nil {
		return nil, err
	}

	var searchCond []string
	//	判断不同数据类型form,table,echarts,bi,3D ....
	children := make([]*RequestParamInfo, 0, len(reqFields))
	for _, field := range reqFields {
		if field.IsSearchCond() {
			searchCond = append(searchCond, field.GetCode())
			continue
		}
		info, err := newRequestParamInfo(field, renderType)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &RequestParams{
		SearchCondList: strings.Join(searchCond, ","),
		RenderType:     renderType,
		Children:       children}, nil
}
