package excelx

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/xuri/excelize/v2"
)

// Marshal 将结构体切片写入Excel文件
func Marshal(filePath string, data interface{}, sheetName ...string) error {
	// 验证输入参数类型
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() != reflect.Slice {
		return errors.New("必须传入切片类型数据")
	}

	// 获取元素类型
	elemType := dataVal.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return errors.New("切片元素必须是结构体")
	}

	// 解析结构体字段信息
	headers, fieldInfos, err := parseStructFields(elemType)
	if err != nil {
		return fmt.Errorf("解析结构体失败: %v", err)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置工作表名称
	sheet := "Sheet1"
	if len(sheetName) > 0 {
		sheet = sheetName[0]
	}

	// 写入表头
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
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
		for colIdx, info := range fieldInfos {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			fieldValue := getFieldStringValue(record.Field(info.index), info.fieldType)
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
func parseStructFields(elemType reflect.Type) ([]string, []structFieldInfo, error) {
	var headers []string
	var fieldInfos []structFieldInfo

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			continue
		}

		headers = append(headers, tag)
		fieldInfos = append(fieldInfos, structFieldInfo{
			index:     i,
			fieldType: field.Type,
		})
	}

	if len(headers) == 0 {
		return nil, nil, errors.New("结构体没有包含excel标签的字段")
	}

	return headers, fieldInfos, nil
}

// 字段信息结构
type structFieldInfo struct {
	index     int
	fieldType reflect.Type
}

// 获取字段字符串值
func getFieldStringValue(field reflect.Value, fieldType reflect.Type) interface{} {
	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		if field.IsNil() {
			return ""
		}
		field = field.Elem()
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

	return fmt.Sprintf("%v", field.Interface())
}
