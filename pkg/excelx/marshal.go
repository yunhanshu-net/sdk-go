package excelx

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/xuri/excelize/v2"
)

// 字段信息结构
type structFieldInfo struct {
	index     int          // 字段索引
	fieldType reflect.Type // 字段类型
}

// Marshal 将结构体切片或指针切片写入Excel文件
// 参数说明：
//   - filePath:  输出的Excel文件路径
//   - data:      要写入的数据，必须是 []Struct 或 []*Struct 类型
//   - sheetName: 可选参数，指定工作表名称（默认Sheet1）
func Marshal(filePath string, data interface{}, sheetName ...string) error {
	// 验证输入数据类型
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() != reflect.Slice {
		return errors.New("必须传入切片类型数据")
	}

	// 获取元素类型并校验
	elemType := dataVal.Type().Elem()
	var structType reflect.Type

	switch {
	case elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct:
		structType = elemType.Elem() // 处理指针类型切片 []*Struct
	case elemType.Kind() == reflect.Struct:
		structType = elemType // 处理值类型切片 []Struct
	default:
		return errors.New("切片元素必须是结构体或结构体指针")
	}

	// 解析结构体字段信息
	headers, fieldInfos, err := parseStructFields(structType)
	if err != nil {
		return fmt.Errorf("结构体解析失败: %v", err)
	}

	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置工作表名称
	sheet := "Sheet1"
	if len(sheetName) > 0 && sheetName[0] != "" {
		sheet = sheetName[0]
	}

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// 写入标题行
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// 写入数据行
	for rowIdx := 0; rowIdx < dataVal.Len(); rowIdx++ {
		record := dataVal.Index(rowIdx)
		var structVal reflect.Value

		// 处理指针类型元素
		if record.Kind() == reflect.Ptr {
			if record.IsNil() {
				continue // 跳过nil指针
			}
			structVal = record.Elem()
		} else {
			structVal = record
		}

		// 写入各字段值
		for colIdx, info := range fieldInfos {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			fieldValue := getFieldStringValue(structVal.Field(info.index), info.fieldType)
			f.SetCellValue(sheet, cell, fieldValue)
		}
	}

	// 自动调整列宽
	for colIdx := range headers {
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		f.SetColWidth(sheet, colName, colName, 18)
	}

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("文件保存失败: %v", err)
	}

	return nil
}

// 解析结构体字段信息
func parseStructFields(structType reflect.Type) ([]string, []structFieldInfo, error) {
	var headers []string
	var fieldInfos []structFieldInfo
	// 遍历所有字段获取excel标签
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			continue // 跳过无标签字段
		}
		headers = append(headers, tag)
		fieldInfos = append(fieldInfos, structFieldInfo{index: i, fieldType: field.Type})
	}
	if len(headers) == 0 {
		return nil, nil, errors.New("结构体中未找到excel标签字段")
	}
	return headers, fieldInfos, nil
}

// 获取字段字符串表示
func getFieldStringValue(field reflect.Value, fieldType reflect.Type) interface{} {
	// 处理指针类型字段
	if fieldType.Kind() == reflect.Ptr {
		if field.IsNil() {
			return "" // nil指针返回空字符串
		}
		field = field.Elem() // 解引用指针
		fieldType = fieldType.Elem()
	}

	// 根据类型转换值
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint()
	case reflect.Float32, reflect.Float64:
		return field.Float()
	case reflect.Bool:
		return field.Bool()
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time)
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// 默认返回字符串表示
	return fmt.Sprintf("%v", field.Interface())
}
