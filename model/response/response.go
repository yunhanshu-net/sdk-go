package response

type Response struct {
	MetaData   map[string]interface{} `json:"meta_data"` //sdk 层
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       *BizData               `json:"body"`
}

type BizData struct {
	MetaData map[string]interface{} `json:"meta_data"` //用户层
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data"`
	Code     int                    `json:"code"`
}

type RunnerResponse struct {
	Response *Response              `json:"response"`
	MetaData map[string]interface{} `json:"meta_data"` //内核层
}

func (r *Response) json(statusCode int, data *BizData) error {
	r.StatusCode = statusCode
	r.Body = data
	return nil
}

func (r *Response) OKWithJSON(data interface{}, meta ...map[string]interface{}) error {
	bz := &BizData{Msg: "ok", Code: 0, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
func (r *Response) FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error {
	bz := &BizData{Msg: msg, Code: -1, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
