package response

type DataType string

const (
	DataTypeJSON = "JSON"
)

type Data struct {
	Type  DataType    `json:"type"`
	Value interface{} `json:"value"`
}
