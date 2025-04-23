package tag

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"reflect"
)

func ParserInputParams(p interface{}) (params []*tagx.FieldInfo, kind reflect.Kind, err error) {
	if p == nil {
		return nil, reflect.Interface, nil
	}
	val := reflect.TypeOf(p)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || val.Kind() != reflect.Slice || val.Kind() != reflect.Array {
		return nil, 0, fmt.Errorf("仅支持返回Struct Slice Array")
	}
	fields, err := tagx.ParseStructFieldsTypeOf(val, "runner")
	if err != nil {
		return nil, 0, err
	}
	return fields, val.Kind(), nil

}
