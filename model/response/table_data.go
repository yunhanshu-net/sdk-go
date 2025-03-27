package response

type Table interface {
	Build() error
}

type tableData struct {
	MetaData   map[string]interface{}   `json:"meta_data"`
	Columns    []column                 `json:"columns"`
	Values     map[string][]interface{} `json:"values"`
	Pagination pagination               `json:"pagination"`
}

func (d *tableData) DataType() DataType {
	return DataTypeTable
}
