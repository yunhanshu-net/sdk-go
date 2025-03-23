package response

type Response struct {
	MetaData   map[string]interface{} `json:"meta_data"` //sdk 层
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	DataType   IDataType              `json:"data_type"`
	Body       interface{}            `json:"body"`
	data       []*IDataType
}

type BizData struct {
	MetaData map[string]interface{} `json:"meta_data"` //用户层
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data"`
	Code     int                    `json:"code"`
}

type RunnerResponse struct {
	Response *Response              `json:"response"`
	MetaData map[string]interface{} `json:"meta_data"` //内核层（runcher引擎层）
}
