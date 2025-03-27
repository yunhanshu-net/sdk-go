package response

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"reflect"
)

func parserTableInfo(row interface{}) []column {
	of := reflect.TypeOf(row)
	var columns []column
	for i := 0; i < of.NumField(); i++ {
		kv := tagx.ParserKv(of.Field(i).Tag.Get("runner"))
		columns = append(columns, column{
			Name: kv["name"],
			Code: kv["code"],
		})
	}
	return columns
}

func (r *Response) Table(dataList interface{}) Table {
	r.DataType = DataTypeTable
	return &table{
		val: dataList,
	}
}

func (t *table) Pagination() {

}

func (t *table) Build() error {
	list, ok := t.val.([]interface{})
	if !ok {
		return fmt.Errorf("类型错误")
	}
	columns := parserTableInfo(list[0])
	marshal, err := sonic.Marshal(t.val)
	if err != nil {
		return err
	}
	var data []map[string]interface{}
	err = sonic.Unmarshal(marshal, &data)
	if err != nil {
		return err
	}
	values := make(map[string][]interface{})
	for _, object := range data {
		for code, value := range object {
			values[code] = append(values[code], value)
		}
	}
	t.columns = columns
	t.values = values

	t.response.DataType = DataTypeTable
	t.response.data = append(t.response.data, &tableData{
		Columns:    t.columns,
		Values:     t.values,
		Pagination: t.pagination,
	})

	return nil
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

type table struct {
	response   *Response
	Error      error
	pagination pagination
	val        interface{}
	columns    []column
	values     map[string][]interface{}
	data       interface{}
}
