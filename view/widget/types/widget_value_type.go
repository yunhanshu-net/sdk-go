package types

var (
	ValueTypes = []string{
		ValueString,
		ValueNumber,
		ValueFloat,
		ValueBoolean,
		ValueArray,
		ValueObject,
		ValueTime,
		ValueFile,
	}
)

func IsValueType(valueType string) bool {
	for _, v := range ValueTypes {
		if valueType == v {
			return true
		}
	}
	return false
}

// 数据类型
const (
	// ValueString 字符串类型
	ValueString = "string"
	// ValueNumber 数字类型
	ValueNumber = "number"
	// ValueBoolean 布尔类型
	ValueBoolean = "boolean"
	// ValueArray 数组类型
	ValueArray = "array"
	// ValueObject 对象类型
	ValueObject = "object"
	// ValueTime 时间类型
	ValueTime = "time"
	// ValueFloat 浮点数类型
	ValueFloat = "float"
	// ValueFile 文件类型
	ValueFile = "file"
)

// UseValueType 这里需要判断，假如请求字段是go的int类型，但是用户不小心在tag里把类型写成string类型了，这时候，我们不应该应用用户的类型，
// 如果应用用户类型导致前端输入字符串到go的int类型，会直接导致请求解析失败，所以这里需要兜底帮用户判断他的类型是否正确，不正确的话我们需要调整正确
func UseValueType(tagType string, fieldType string) string {
	if tagType != fieldType {
		return tagType
	}
	return fieldType //用field的类型是安全的
}
