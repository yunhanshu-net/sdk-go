package response

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"gorm.io/gorm"
	"net/http"
	"reflect"
)

type Table interface {
	Builder
	AutoPaginated(dbAndWhere *gorm.DB, model interface{}, pageInfo *request.PageInfo) Table
}

// Paginated 分页查询结果结构体
type Paginated struct {
	CurrentPage int `json:"current_page"` // 当前页码
	TotalCount  int `json:"total_count"`  // 总数据量
	TotalPages  int `json:"total_pages"`  // 总页数
	PageSize    int `json:"page_size"`    // 每页数量
}

func (r *Data) Table(resultList interface{}, title ...string) Table {
	titleStr := ""
	if len(title) > 0 {
		titleStr = title[0]
	}
	r.StatusCode = http.StatusOK
	return &tableData{
		val:        resultList,
		response:   r,
		title:      titleStr,
		pagination: &Paginated{},
	}
}

type column struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}
type tableData struct {
	TraceID    string `json:"trace_id"`
	error      error
	title      string
	response   *Data
	pagination *Paginated
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

func (t *tableData) AutoPaginated(db *gorm.DB, model interface{}, pageInfo *request.PageInfo) Table {
	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := db.Model(model).Count(&totalCount).Error; err != nil {
		t.error = fmt.Errorf("AutoPaginated.Count :%+v failed to count records: %v", t.val, err)
		return t
	}

	if pageInfo.GetSorts() != "" {
		db.Order(pageInfo.GetSorts())
	}
	// 查询当前页数据
	if err := db.Offset(offset).Limit(pageSize).Find(t.val).Error; err != nil {
		t.error = fmt.Errorf("AutoPaginated.Find :%+v failed to count records: %v", t.val, err)
		return t
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	t.pagination = &Paginated{
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}
	return t
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
		Pagination Paginated                `json:"pagination"`
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
