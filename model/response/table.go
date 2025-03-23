package response

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func parserKv(tag string) map[string]string {
	mp := make(map[string]string)
	split := strings.Split(tag, ";")
	for _, s := range split {
		vals := strings.Split(s, ":")
		key := vals[0]
		value := vals[1]
		mp[key] = value
	}
	return mp
}
func parserTableInfo(row interface{}) []column {
	of := reflect.TypeOf(row)
	var columns []column
	for i := 0; i < of.NumField(); i++ {
		kv := parserKv(of.Field(i).Tag.Get("runner"))
		columns = append(columns, column{
			Name: kv["name"],
			Code: kv["code"],
		})
	}
	return columns
}

func (r *Response) Table(dataList interface{}) *Table {
	return &Table{
		val: dataList,
	}
}

func (t *Table) Pagination() {

}

func (t *Table) Build() *Table {
	list, ok := t.val.([]interface{})
	if !ok {
		return &Table{Error: fmt.Errorf("类型错误")}
	}
	columns := parserTableInfo(list[0])
	marshal, err := json.Marshal(t.val)
	if err != nil {
		return &Table{Error: err}
	}
	var data []map[string]interface{}
	err = json.Unmarshal(marshal, &data)
	if err != nil {
		return &Table{Error: err}
	}
	values := make(map[string][]interface{})
	for _, object := range data {
		for code, value := range object {
			values[code] = append(values[code], value)
		}
	}
	return &Table{
		columns: columns,
		values:  values,
	}
}

type column struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type pagination struct {
	currentPage int
	totalPage   int
	pageSize    int
	totalCount  int
}

type Table struct {
	Error      error
	pagination pagination
	val        interface{}
	columns    []column
	values     map[string][]interface{}
	data       interface{}
}

func (t *Table) name() {

}
