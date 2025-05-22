package api

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"reflect"
)

func GetFields(el interface{}) ([]*tagx.RunnerFieldInfo, error) {
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
	return tags, nil
}
