package tagx

import (
	"fmt"
	"reflect"
	"strings"
)

func ParserKv(tag string) map[string]string {
	if tag == "" {
		return make(map[string]string)
	}
	mp := make(map[string]string)
	split := strings.Split(tag, ";")
	for _, s := range split {
		vals := strings.Split(s, ":")
		key := vals[0]
		value := vals[1]
		mp[key] = value
	}
	return mp
}

type FieldInfo struct {
	Name  string            // 字段名（含匿名字段层级，如"User.ID"）
	Type  reflect.Type      // 字段类型
	Tags  map[string]string // 解析后的标签键值对
	Index []int             // 字段索引路径（用于反射）
}

func ParseStructFields(obj interface{}, tagKey string) ([]FieldInfo, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}
	return parseFields(t, tagKey, nil, nil), nil
}

// 递归解析字段，处理匿名字段
func parseFields(t reflect.Type, tagKey string, parentIndex []int, parentNames []string) []FieldInfo {
	var fields []FieldInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		currentIndex := append(parentIndex, i)
		currentNames := append(parentNames, field.Name)

		// 处理匿名字段（结构体类型）
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedFields := parseFields(field.Type, tagKey, currentIndex, parentNames)
			fields = append(fields, embeddedFields...)
			continue
		}

		// 普通字段：生成FieldInfo
		info := FieldInfo{
			Name:  strings.Join(currentNames, "."), // 生成层级字段名
			Type:  field.Type,
			Tags:  parseTag(field.Tag.Get(tagKey)),
			Index: currentIndex,
		}
		fields = append(fields, info)
	}
	return fields
}

// 解析标签字符串为键值对（支持GORM风格分号分隔）
func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(tag, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, ":", 2)
		key := strings.TrimSpace(kv[0])
		if key == "" {
			continue
		}
		value := ""
		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
		}
		result[key] = value
	}
	return result
}
