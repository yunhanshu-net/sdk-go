package request

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Paginated struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"current_page"` // 当前页码
	TotalCount  int         `json:"total_count"`  // 总数据量
	TotalPages  int         `json:"total_pages"`  // 总页数
	PageSize    int         `json:"page_size"`    // 每页数量
}

func AutoPaginated(db *gorm.DB, model interface{}, data interface{}, pageInfo *PageInfo) (*Paginated, error) {

	if pageInfo == nil {
		pageInfo = new(PageInfo)
	}
	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := db.Model(model).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("AutoPaginated.Count :%+v failed to count records: %v", data, err)

	}

	if pageInfo.GetSorts() != "" {
		db.Order(pageInfo.GetSorts())
	}
	// 查询当前页数据
	if err := db.Offset(offset).Limit(pageSize).Find(data).Error; err != nil {
		return nil, fmt.Errorf("AutoPaginated.Find :%+v failed to count records: %v", data, err)
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	return &Paginated{
		Data:        data,
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}, nil

}

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`

	//age,desc,score,asc,height,desc
	Sorts string `json:"sorts" form:"sorts"` //
}

// GetLimit 获取分页大小，支持默认值
func (i *PageInfo) GetLimit(defaultSize ...int) int {
	if i.PageSize <= 0 { // 如果 PageSize 小于等于 0
		if len(defaultSize) > 0 {
			return defaultSize[0] // 使用传入的默认值
		}
		return 20 // 使用固定默认值 20
	}
	return i.PageSize // 返回 PageSize
}

// GetOffset 获取分页偏移量
func (i *PageInfo) GetOffset() int {
	if i.Page < 1 { // 如果 Page 小于 1，设置为 1
		i.Page = 1
	}
	return (i.Page - 1) * i.GetLimit() // 计算偏移量
}

// 定义排序字段结构体
type sortField struct {
	Field string // 字段名
	Order string // 排序方向，"asc" 或 "desc"
}

// 解析排序字段字符串
func parseSortFields(sortStr string) ([]sortField, error) {
	if sortStr == "" {
		return nil, errors.New("排序字段不能为空")
	}

	// 按逗号分割字符串
	parts := strings.Split(sortStr, ",")
	if len(parts)%2 != 0 {
		return nil, errors.New("排序字段格式错误：字段名和排序方向必须成对出现")
	}

	var sortFields []sortField
	for i := 0; i < len(parts); i += 2 {
		field := strings.TrimSpace(parts[i])   // 字段名
		order := strings.TrimSpace(parts[i+1]) // 排序方向

		order = strings.ToUpper(strings.TrimSpace(order))
		// 校验排序方向
		if order != "ASC" && order != "DESC" {
			return nil, fmt.Errorf("无效的排序方向：%s", order)
		}

		// 添加到结果切片
		sortFields = append(sortFields, sortField{Field: field, Order: order})
	}

	return sortFields, nil
}

func (i *PageInfo) GetSorts() string {
	sortFields, err := parseSortFields(i.Sorts)
	if err != nil {
		return ""
	}
	str := &strings.Builder{}
	for idx, s := range sortFields {
		str.WriteString(s.Field) //这里做好sql的防注入
		str.WriteString(s.Order)
		if idx != len(sortFields)-1 {
			str.WriteString(",")
		}
	}
	return str.String()
}
