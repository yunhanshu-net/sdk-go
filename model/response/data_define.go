package response

type DataType string

const (
	DataTypeJSON    = "json"
	DataTypeTable   = "table"
	DataTypeFiles   = "files"
	DataTypeEcharts = "echarts"
)

type Data struct {
	Type  DataType    `json:"type"`
	Value interface{} `json:"value"`
}
