package api

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"reflect"
)

type TableRequestParams struct {
	RenderType string                  `json:"render_type"`
	Children   []*FormRequestParamInfo `json:"children"`
}

func NewTableRequestParams(el interface{}) (*TableRequestParams, error) {

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
	children := make([]*FormRequestParamInfo, 0, len(reqFields))
	for _, field := range reqFields {
		if field.IsSearchCond() {
			searchCond = append(searchCond, field.GetCode())
			continue
		}
		info, err := newFormRequestParamInfo(field, response.RenderTypeTable)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}
	return &TableRequestParams{
		RenderType: response.RenderTypeTable,
		Children:   children,
	}, nil
}
