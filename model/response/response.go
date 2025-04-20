package response

type Response interface {
	JSON(data interface{}) JSON
	FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error
	Table(dataList interface{}, desc ...string) Table
}

type Data struct {
	MetaData   map[string]interface{} `json:"meta_data"`   //sdk 层
	StatusCode int                    `json:"status_code"` //http对应http code 正常200
	Msg        string                 `json:"msg"`
	Headers    map[string]string      `json:"headers"`
	DataType   DataType               `json:"data_type"`
	Body       interface{}            `json:"body"`
	Multiple   bool                   `json:"multiple"` //是否是多个数据类型？比如返回
	data       anyData
}

func (d *Data) SetMetaData(key string, value interface{}) {
	if d.MetaData == nil {
		d.MetaData = make(map[string]interface{})
	}
	d.MetaData[key] = value
}

type BizData struct {
	MetaData map[string]interface{} `json:"meta_data"` //用户层
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data_list"`
	Code     int                    `json:"code"`
}

type RunnerResponse struct {
	Response *Data                  `json:"response"`
	MetaData map[string]interface{} `json:"meta_data"` //内核层（runcher引擎层）
}
