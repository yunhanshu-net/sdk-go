package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/common/slicesx"
	"github.com/yunhanshu-net/sdk-go/pkg/stringsx"
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

	return &Params{RenderType: stringsx.DefaultString(t, response.RenderTypeForm), Children: children}, nil
}

func NewResponseParams(el interface{}, renderType string) (*Params, error) {

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
		paramsOut, err := newParams(resFields, renderType)
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
		paramsOut, err := newParams(of, renderType)
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
	validate := tag.Type.Tag.Get("validate")
	split := strings.Split(validate, ",")

	param := &ParamInfo{
		Code:         tag.Tags["code"],
		Name:         tag.Tags["name"],
		Desc:         tag.Tags["desc"],
		Required:     slicesx.ContainsString(split, "required"),
		Validates:    strings.Join(slicesx.RemoveString(split, "required"), ","),
		Callbacks:    tag.Tags["callback"],
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

func newParams(fields []*tagx.FieldInfo, renderType string) (*Params, error) {
	//	判断不同数据类型form,table,echarts,bi,3D ....
	children := make([]*ParamInfo, 0, len(fields))
	for _, field := range fields {
		info, err := newParamInfo(field)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &Params{RenderType: stringsx.DefaultString(renderType, response.RenderTypeForm), Children: children}, nil
}
