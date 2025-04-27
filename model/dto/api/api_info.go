package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	render2 "github.com/yunhanshu-net/sdk-go/view/render"
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
	p := &ParamInfo{
		Code: tag.Tags["code"],
		Name: tag.Tags["name"],
		Desc: tag.Tags["desc"],
	}

	required, ok := tag.Tags["required"]
	if ok {
		if required == "" {
			p.Required = true
		} else {
			p.Required = required == "true"
		}
		return p, nil
	}

	if p.Code == "" {
		get := tag.Type.Tag.Get("json")
		if get != "" {
			p.Code = strings.Split(get, ",")[0]
		}
	}

	if p.Name == "" {
		p.Name = p.Code
	}
	widget, err := render2.NewWidget(tag)
	if err != nil {
		return nil, err
	}
	p.Widget = widget
	return p, nil
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
	Code string `json:"code,omitempty"`
	//中文名称
	Name string `json:"name,omitempty"`
	//中文介绍
	Desc string `json:"desc,omitempty"`
	//是否必填
	Required bool `json:"required,omitempty"`

	Widget render2.Widget `json:"widget"`
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
	Callbacks []string `json:"callbacks"`
}

type ApiLogs struct {
	Version string  `json:"version"`
	Apis    []*Info `json:"apis"`
}
