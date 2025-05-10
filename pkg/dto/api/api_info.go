package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"github.com/yunhanshu-net/sdk-go/view/widget"
	"github.com/yunhanshu-net/sdk-go/view/widget/types"
	"reflect"
	"strings"
)

func NewRequestParams(el interface{}, t string) (*Params, error) {
	typeOf := reflect.TypeOf(el)
	logrus.Info("typeOf: kind", typeOf.Kind())
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

	//	判断不同数据类型form,table,echarts,bi,3D ....
	children := make([]*ParamInfo, 0, len(reqFields))
	for _, field := range reqFields {
		info, err := newParamInfo(field)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &Params{RenderType: t, Children: children}, nil
}

func NewResponseParams(el interface{}, t string) (*Params, error) {

	rspType := reflect.TypeOf(el)
	logrus.Info("rspType: kind", rspType.Kind())
	if rspType.Kind() == reflect.Pointer {
		rspType = rspType.Elem()
	}
	if rspType.Kind() != reflect.Struct && rspType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("输出参数仅支持Struct和Slice类型")
	}

	if rspType.Kind() == reflect.Struct {
		resFields, err := tagx.ParseStructFieldsTypeOf(rspType, "runner")
		if err != nil {
			return nil, err
		}
		paramsOut, err := newParams(resFields, t)
		if err != nil {
			return nil, err
		}
		return paramsOut, nil
	} else {
		tp, err := tagx.GetSliceElementType(el)
		if err != nil {
			return nil, err
		}
		of, err := tagx.ParseStructFieldsTypeOf(tp, "runner")
		if err != nil {
			return nil, err
		}
		paramsOut, err := newParams(of, t)
		if err != nil {
			return nil, err
		}
		return paramsOut, nil
	}
}
func newParamInfo(tag *tagx.FieldInfo) (*ParamInfo, error) {
	if tag == nil {
		return nil, fmt.Errorf("tag==nil")
	}
	if tag.Tags == nil {
		tag.Tags = map[string]string{}
	}
	widgetIns, err := widget.NewWidget(tag)
	if err != nil {
		return nil, err
	}
	valueType, err := tag.GetValueType()
	if err != nil {
		return nil, err
	}
	if !types.IsValueType(valueType) {
		return nil, fmt.Errorf("不是合法的值类型：%s", valueType)
	}

	param := &ParamInfo{
		Code:         tag.Tags["code"],
		Name:         tag.Tags["name"],
		Desc:         tag.Tags["desc"],
		Validates:    tag.Type.Tag.Get("validate"),
		Callbacks:    tag.Type.Tag.Get("callback"),
		WidgetConfig: widgetIns,
		WidgetType:   widgetIns.GetWidgetType(),
		ValueType:    types.UseValueType(tag.Tags["type"], valueType),
		Example:      tag.Tags["example"],
	}

	if param.Code == "" {
		get := tag.Type.Tag.Get("json")
		if get != "" {
			param.Code = strings.Split(get, ",")[0]
		}
	}

	if param.Name == "" {
		param.Name = param.Code
	}

	return param, nil
}

func newParams(fields []*tagx.FieldInfo, t string) (*Params, error) {
	//	判断不同数据类型form,table,echarts,bi,3D ....
	children := make([]*ParamInfo, 0, len(fields))
	for _, field := range fields {
		info, err := newParamInfo(field)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &Params{RenderType: t, Children: children}, nil
}

type Params struct {

	//form,table,echarts,bi,3D .....
	RenderType string       `json:"render_type"`
	Children   []*ParamInfo `json:"children"`
}

type ParamInfo struct {
	//英文标识
	Code string `json:"code"`
	//中文名称
	Name string `json:"name"`
	//中文介绍
	Desc string `json:"desc"`
	//是否必填
	Required bool `json:"required"`

	Callbacks    string        `json:"callbacks"`
	Validates    string        `json:"validates"`
	WidgetConfig widget.Widget `json:"widget_config"`
	WidgetType   string        `json:"widget_type"`
	ValueType    string        `json:"value_type"`
	Example      string        `json:"example"`
}

type Info struct {
	Router      string   `json:"router"`
	Method      string   `json:"method"`
	User        string   `json:"user"`
	Runner      string   `json:"runner"`
	ApiDesc     string   `json:"api_desc"`
	Labels      []string `json:"labels"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	Classify    string   `json:"classify"`
	//输入参数
	ParamsIn *Params `json:"params_in"`
	//输出参数
	ParamsOut *Params  `json:"params_out"`
	UseTables []string `json:"use_tables"`
	UseDB     []string `json:"use_db"`
	Callbacks []string `json:"callbacks"`
}

type ApiLogs struct {
	Version string  `json:"version"`
	Apis    []*Info `json:"apis"`
}
