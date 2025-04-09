package response

type rsp struct {
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	DataType DataType               `json:"data_type"`
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Data     interface{}            `json:"data"`
}
