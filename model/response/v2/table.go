package v2

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"net/http"
	"reflect"
)

type Table interface {
	Builder
}

func (r *ResponseData) Table(dataList interface{}, title ...string) Table {
	titleStr := ""
	if len(title) > 0 {
		titleStr = title[0]
	}
	r.StatusCode = http.StatusOK
	return &tableData{
		val:        dataList,
		response:   r,
		title:      titleStr,
		pagination: &pagination{},
	}
}

type column struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}
type pagination struct {
	currentPage int
	totalPage   int
	pageSize    int
	totalCount  int
}
type tableData struct {
	TraceID    string `json:"trace_id"`
	title      string
	response   *ResponseData
	pagination *pagination
	val        interface{}
	columns    []column
	values     map[string][]interface{}
	data       interface{}

	buildData string
}

func (t *tableData) Build() error {
	sliceVal := reflect.ValueOf(t.val)
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("类型错误")
	}
	if sliceVal.Len() == 0 {
		return fmt.Errorf("列表为空")
	}
	row := sliceVal.Index(0)
	columns := parserTableInfo(row.Interface())
	//marshal, err := sonic.Marshal(t.val)
	//if err != nil {
	//	return err
	//}
	//var data []map[string]interface{}
	//err = sonic.Unmarshal(marshal, &data)
	//if err != nil {
	//	return err
	//}
	//values := make(map[string][]interface{})
	//for _, object := range data {
	//	for code, value := range object {
	//		values[code] = append(values[code], value)
	//	}
	//}

	values := make(map[string][]interface{}, len(columns))
	for i := 0; i < sliceVal.Len(); i++ {
		row := sliceVal.Index(i)
		for _, col := range columns {
			field := row.Field(col.Idx) //根据idx取field
			if _, ok := values[col.Code]; !ok {
				values[col.Code] = make([]interface{}, 0, sliceVal.Len())
			}
			values[col.Code] = append(values[col.Code], field.Interface())
		}
	}

	t.columns = columns
	t.values = values
	t.response.DataType = DataTypeTable
	err := t.buildJSON()
	if err != nil {
		return err
	}
	t.response.Multiple = false
	t.response.DataType = t.DataType()
	return nil
}

func (t *tableData) buildJSON() error {
	type rsp struct {
		Code     int         `json:"code"`
		Msg      string      `json:"msg"`
		DataType DataType    `json:"data_type"`
		Data     interface{} `json:"data"`
	}

	type Table struct {
		Title      string                   `json:"title"`
		Column     []column                 `json:"column"`
		Values     map[string][]interface{} `json:"values"`
		Pagination pagination               `json:"pagination"`
	}
	//Title      string                   `json:"title"`
	//Column     []column                 `json:"column"`
	//Values     map[string][]interface{} `json:"values"`
	//Pagination pagination               `json:"pagination"`

	tb := rsp{
		Code:     successCode,
		Msg:      successMsg,
		DataType: DataTypeTable,
		Data: Table{
			Title:      t.title,
			Column:     t.columns,
			Values:     t.values,
			Pagination: *t.pagination,
		},
	}
	t.response.Body = tb
	return nil
}

func (t *tableData) BuildJSON() string {
	t.response.Multiple = false
	t.response.DataType = t.DataType()
	t.response.Body = t.buildData
	return t.buildData
}

func (t *tableData) DataType() DataType {
	return DataTypeTable
}

func parserTableInfo(row interface{}) []column {
	of := reflect.TypeOf(row)
	var columns []column
	for i := 0; i < of.NumField(); i++ {
		kv := tagx.ParserKv(of.Field(i).Tag.Get("runner"))
		columns = append(columns, column{Name: kv["name"], Code: kv["code"], Idx: i})
	}
	return columns
}
