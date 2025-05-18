package response

type RunFunctionResp struct {
	MetaData map[string]interface{} `json:"meta_data"`
	Headers  map[string]string      `json:"headers"`
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	TraceID  string                 `json:"trace_id"`
	DataType string                 `json:"data_type"`
	Data     interface{}            `json:"data"`
	DataList []interface{}          `json:"data_list"`
	Multiple bool                   `json:"multiple"`
}

type Builder interface {
	Build() error
}

type Form interface {
	Builder
}

type Response interface {
	Form(data interface{}) Form
}

func (r *RunFunctionResp) Form(data interface{}) Form {
	return newForm(data, r)
}
