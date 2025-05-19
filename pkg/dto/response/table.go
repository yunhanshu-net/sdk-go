package response

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/query"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"gorm.io/gorm"
	"reflect"
)

type Table interface {
	Builder
	AutoPaginated(dbAndWhere *gorm.DB, model interface{}, pageInfo *query.PageInfoReq) Table
}

type column struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}
type paginated struct {
	CurrentPage int `json:"current_page"` // 当前页码
	TotalCount  int `json:"total_count"`  // 总数据量
	TotalPages  int `json:"total_pages"`  // 总页数
	PageSize    int `json:"page_size"`    // 每页数量
}

type tableData struct {
	err  error
	val  interface{}
	resp *RunFunctionResp
	Data table
}
type table struct {
	Title      string                   `json:"title"`
	Column     []column                 `json:"column"`
	Values     map[string][]interface{} `json:"values"`
	Pagination paginated                `json:"pagination"`
}

func newTable(resp *RunFunctionResp, resultList interface{}, title ...string) *tableData {
	titleStr := ""
	if len(title) > 0 {
		titleStr = title[0]
	}
	return &tableData{resp: resp, val: resultList, Data: table{Title: titleStr}}
}
func (r *RunFunctionResp) Table(resultList interface{}, title ...string) Table {
	return newTable(r, resultList, title...)
}

func (t *tableData) Build() error {
	if t.err != nil {
		return t.err
	}
	sliceVal := reflect.ValueOf(t.val)
	if sliceVal.Kind() == reflect.Pointer {
		sliceVal = sliceVal.Elem()
	}
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

	t.Data.Column = columns
	t.Data.Values = values
	return build(t.resp, t.Data, RenderTypeTable)
}

func (t *tableData) AutoPaginated(db *gorm.DB, model interface{}, pageInfo *query.PageInfoReq) Table {

	if pageInfo == nil {
		pageInfo = new(query.PageInfoReq)
	}
	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := db.Model(model).Count(&totalCount).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Count :%+v failed to count records: %v", t.val, err)
		return t
	}

	if pageInfo.GetSorts() != "" {
		db.Order(pageInfo.GetSorts())
	}
	// 查询当前页数据
	if err := db.Offset(offset).Limit(pageSize).Find(t.val).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Find :%+v failed to count records: %v", t.val, err)
		return t
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	t.Data.Pagination = paginated{
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}
	return t
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
