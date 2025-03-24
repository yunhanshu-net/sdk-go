// Package jsonx ...
package jsonx

import (
	"fmt"
	"github.com/bytedance/sonic"
)

// MustJSON ..
func MustJSON(el interface{}) string {
	marshal, err := sonic.Marshal(el)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// MustPrintJSON ...
func MustPrintJSON(el interface{}) {
	marshal, err := sonic.Marshal(el)
	if err != nil {
		fmt.Println(fmt.Sprintf("[jsonx] err:%s el:%+v", err.Error(), el))
		return
	}
	fmt.Println(string(marshal))
}

// JSONString ...
func JSONString(el interface{}) string {
	marshal, err := sonic.Marshal(el)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// String ...
func String(el interface{}) string {
	marshal, err := sonic.Marshal(el)
	if err != nil {
		return ""
	}
	return string(marshal)
}
