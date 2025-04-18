package excelx

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
	"time"
)

// Unmarshal 将Excel文件内容解析到结构体切片中
func Unmarshal(filePath string, out interface{}, sheetName ...string) error {
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("必须传入结构体切片指针")
	}

	sliceVal := val.Elem()
	elemType := sliceVal.Type().Elem()

	//if elemType.Kind() != reflect.Struct {
	//	//[]User || []*User
	//	return errors.New("目标必须是结构体切片")
	//}

	//if !(elemType.Kind() == reflect.Struct ||
	//	(elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct)) {
	//	return errors.New("目标必须是结构体切片或结构体指针切片")
	//}

	var structType reflect.Type
	switch {
	case elemType.Kind() == reflect.Struct:
		structType = elemType
	case elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct:
		structType = elemType.Elem()
	default:
		return errors.New("目标必须是结构体或结构体指针切片")
	}

	// 解析结构体标签
	fieldMap := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if tag := field.Tag.Get("excel"); tag != "" {
			fieldMap[tag] = i
		}
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("文件打开失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return errors.New("excel文件没有工作表")
	}
	name := sheets[0]
	if len(sheetName) > 0 {
		exist := false
		for _, s := range sheets {
			if s == sheetName[0] {
				exist = true
				name = s
				break
			}
		}
		if !exist {
			return fmt.Errorf("sheetName %s not exist", name)
		}
	}

	rows, err := f.GetRows(name)
	if err != nil {
		return fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) == 0 {
		return nil
	}

	// 解析标题行
	headers := rows[0]
	colIndexMap := make(map[string]int)
	for idx, header := range headers {
		colIndexMap[header] = idx
	}

	// 处理数据行（关键修改点）
	for rowIdx, row := range rows[1:] {
		//newElem := reflect.New(elemType).Elem()

		var newElem reflect.Value
		if elemType.Kind() == reflect.Ptr {
			// 指针类型：创建结构体实例并取其地址
			structType := elemType.Elem()
			newElem = reflect.New(structType) // 生成 *structType
		} else {
			// 结构体类型：直接创建实例
			newElem = reflect.New(elemType).Elem() // 生成 structType
		}

		for tag, fieldIndex := range fieldMap {
			// 检查列是否存在
			colIndex, exists := colIndexMap[tag]
			if !exists {
				continue // 跳过不存在的列
			}

			// 处理可能为空的值
			var cellValue string
			if colIndex < len(row) {
				cellValue = row[colIndex]
			}

			// 设置字段值
			//field := newElem.Field(fieldIndex)
			var structVal reflect.Value
			if elemType.Kind() == reflect.Ptr {
				structVal = newElem.Elem() // 解引用指针
			} else {
				structVal = newElem
			}
			field := structVal.Field(fieldIndex)

			if err := setValue(field, cellValue); err != nil {
				return fmt.Errorf("第%d行[%s]解析错误: %v", rowIdx+2, tag, err)
			}
		}

		sliceVal.Set(reflect.Append(sliceVal, newElem))
	}

	return nil
}

// 修改后的setValue函数
func setValue(field reflect.Value, value string) error {
	// 处理指针类型
	if field.Kind() == reflect.Ptr {
		return handlePointerType(field, value)
	}

	// 处理空值（非指针类型）
	if value == "" {
		return setZeroValue(field)
	}

	// 基础类型处理
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			return parseTime(field, value)
		}
		return fmt.Errorf("不支持的struct类型: %s", field.Type())
	default:
		return fmt.Errorf("不支持的字段类型: %s", field.Type())
	}
	return nil
}

// 处理指针类型的专用函数
func handlePointerType(field reflect.Value, value string) error {
	// 如果值为空，直接设置为nil
	if value == "" {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	// 创建指针指向的新实例
	elemType := field.Type().Elem()
	newValue := reflect.New(elemType)

	// 递归设置值
	if err := setValue(newValue.Elem(), value); err != nil {
		return err
	}

	field.Set(newValue)
	return nil
}

// 设置非指针类型的零值
func setZeroValue(field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(0)
	case reflect.Float32, reflect.Float64:
		field.SetFloat(0)
	case reflect.Bool:
		field.SetBool(false)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			field.Set(reflect.ValueOf(time.Time{}))
			return nil
		}
		return fmt.Errorf("不支持的struct类型: %s", field.Type())
	default:
		return fmt.Errorf("不支持的字段类型: %s", field.Type())
	}
	return nil
}

// 时间解析示例（可根据需要扩展格式）
func parseTime(field reflect.Value, value string) error {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	for _, format := range formats {
		t, err := time.Parse(format, value)
		if err == nil {
			field.Set(reflect.ValueOf(t))
			return nil
		}
	}
	return fmt.Errorf("无法解析时间格式: %s", value)
}
