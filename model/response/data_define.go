package response

type DataType string

const (
	DataTypeForm    = "form"
	DataTypeJSON    = "json"
	DataTypeTable   = "table"
	DataTypeFiles   = "files"
	DataTypeEcharts = "echarts"
)

type Builder interface {
	Build() error
}

type anyData interface {
	Builder
	BuildJSON() string
	GetDataType() DataType
}
