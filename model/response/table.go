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

// paginated 分页查询结果结构体
type paginated struct {
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
		TraceID:    r.TraceID,
		val:        resultList,
		response:   r,
		title:      titleStr,
		pagination: &paginated{},
	}
}

type column struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type tableResp struct {
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	DataType RenderType             `json:"data_type"`
	Data     interface{}            `json:"data"`
}

type table struct {
	Title      string                   `json:"title"`
	Column     []column                 `json:"column"`
	Values     map[string][]interface{} `json:"values"`
	Pagination paginated                `json:"pagination"`
}

type tableData struct {
	TraceID    string `json:"trace_id"`
	error      error
	title      string
	response   *Data
	pagination *paginated
	val        interface{}
	columns    []column
	values     map[string][]interface{}
	data       interface{}

	buildData string
}

func (t *tableData) Build() error {
	sliceVal := reflect.ValueOf(t.val)
	fmt.Println("类型", sliceVal.Kind())
	if sliceVal.Kind() == reflect.Pointer {
		sliceVal = sliceVal.Elem()
	}
	fmt.Println("类型", sliceVal.Kind())
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
	t.response.RenderType = RenderTypeTable
	err := t.buildJSON()
	if err != nil {
		return err
	}
	t.response.Multiple = false
	t.response.RenderType = t.DataType()
	return nil
}

func (t *tableData) AutoPaginated(db *gorm.DB, model interface{}, pageInfo *request.PageInfo) Table {

	if pageInfo == nil {
		pageInfo = new(request.PageInfo)
	}
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
	t.pagination = &paginated{
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}
	return t
}

func (t *tableData) buildJSON() error {
	tb := tableResp{
		TraceID:  t.TraceID,
		MetaData: nil,
		Code:     successCode,
		Msg:      successMsg,
		DataType: RenderTypeTable,
		Data: table{
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
	t.response.RenderType = t.DataType()
	t.response.Body = t.buildData
	return t.buildData
}

func (t *tableData) DataType() RenderType {
	return RenderTypeTable
}

func parserTableInfo(row interface{}) []column {
	of := reflect.TypeOf(row)
	var columns []column
	for i := 0; i < of.NumField(); i++ {
		field := of.Field(i)
		kv := tagx.ParserKv(field.Tag.Get("runner"))
		name := kv["name"]
		code := kv["code"]
		if name == "" {
			name = field.Tag.Get("json")
			if name == "" {
				name = field.Name
			}
		}
		if code == "" {
			name = field.Tag.Get("json")
			if name == "" {
				name = field.Name
			}
		}

		columns = append(columns, column{Name: name, Code: code, Idx: i})
	}
	return columns
}
