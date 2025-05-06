package base

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// Paginated 分页结果结构体
type Paginated[T any] struct {
	Items       T     `json:"items"`        // 分页数据
	CurrentPage int   `json:"current_page"` // 当前页码
	TotalCount  int64 `json:"total_count"`  // 总数据量
	TotalPages  int   `json:"total_pages"`  // 总页数
	PageSize    int   `json:"page_size"`    // 每页数量
}

// PageInfoReq 分页参数结构体
type PageInfoReq struct {
	Page     int    `json:"page" form:"page" binding:"omitempty,min=1"`           // 页码，从1开始
	PageSize int    `json:"page_size" form:"page_size" binding:"omitempty,min=1"` // 每页记录数
	Sorts    string `json:"sorts" form:"sorts"`                                   // 排序字段，格式：field1,asc,field2,desc
}

// sortField 排序字段结构体
type sortField struct {
	Field string // 字段名
	Order string // 排序方向，"ASC" 或 "DESC"
}

// GetLimit 获取分页大小，支持默认值
func (i *PageInfoReq) GetLimit(defaultSize ...int) int {
	if i.PageSize <= 0 { // 如果 PageSize 小于等于 0
		if len(defaultSize) > 0 {
			return defaultSize[0] // 使用传入的默认值
		}
		return 20 // 使用固定默认值 20
	}
	return i.PageSize // 返回 PageSize
}

// GetOffset 获取分页偏移量
func (i *PageInfoReq) GetOffset() int {
	if i.Page < 1 { // 如果 Page 小于 1，设置为 1
		i.Page = 1
	}
	return (i.Page - 1) * i.GetLimit() // 计算偏移量
}

// SafeColumn 检查列名是否安全（防SQL注入）
func SafeColumn(column string) bool {
	// 列名只允许字母、数字、下划线
	for _, c := range column {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// ParseSortFields 解析排序字段字符串
func ParseSortFields(sortStr string) ([]sortField, error) {
	if sortStr == "" {
		return nil, nil
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

		// 检查字段名是否安全（防止SQL注入）
		if !SafeColumn(field) {
			return nil, fmt.Errorf("无效的排序字段名：%s", field)
		}

		order = strings.ToUpper(order)
		// 校验排序方向
		if order != "ASC" && order != "DESC" {
			return nil, fmt.Errorf("无效的排序方向：%s", order)
		}

		// 添加到结果切片
		sortFields = append(sortFields, sortField{Field: field, Order: order})
	}

	return sortFields, nil
}

// GetSorts 获取排序SQL
func (i *PageInfoReq) GetSorts() string {
	sortFields, err := ParseSortFields(i.Sorts)
	if err != nil || len(sortFields) == 0 {
		return ""
	}

	var orderClauses []string
	for _, s := range sortFields {
		// 再次检查字段名是否安全
		if !SafeColumn(s.Field) {
			continue
		}
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", s.Field, s.Order))
	}

	return strings.Join(orderClauses, ", ")
}

// AutoPaginate 自动分页查询
func AutoPaginate[T any](ctx context.Context, db *gorm.DB, model interface{}, data T, pageInfo *PageInfoReq) (*Paginated[T], error) {
	if pageInfo == nil {
		pageInfo = new(PageInfoReq)
	}

	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := db.Model(model).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("分页查询统计总数失败: %w", err)
	}

	// 应用排序条件
	sortStr := pageInfo.GetSorts()
	if sortStr != "" {
		db = db.Order(sortStr)
	}

	// 查询当前页数据
	if err := db.Offset(offset).Limit(pageSize).Find(data).Error; err != nil {
		return nil, fmt.Errorf("分页查询数据失败: %w", err)
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	return &Paginated[T]{
		Items:       data,
		CurrentPage: pageInfo.Page,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		PageSize:    pageSize}, nil
}
