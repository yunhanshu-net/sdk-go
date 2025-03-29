package v2

type Response interface {
	JSON(data interface{}) JSON
	FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error
	Table(dataList interface{}, desc ...string) Table
}

type ResponseData struct {
	MetaData   map[string]interface{} `json:"meta_data"`   //sdk 层
	StatusCode int                    `json:"status_code"` //http对应http code 正常200
	Headers    map[string]string      `json:"headers"`
	DataType   DataType               `json:"data_type"`
	Body       interface{}            `json:"body"`
	Multiple   bool                   `json:"multiple"` //是否是多个数据类型？比如返回
	//dataList   []anyData
	data anyData
}

type BizData struct {
	MetaData map[string]interface{} `json:"meta_data"` //用户层
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data_list"`
	Code     int                    `json:"code"`
}

type RunnerResponse struct {
	Response *ResponseData          `json:"response"`
	MetaData map[string]interface{} `json:"meta_data"` //内核层（runcher引擎层）
}
