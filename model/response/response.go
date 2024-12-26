package response

import "encoding/json"

type BizData struct {
	MetaData map[string]interface{} `json:"meta_data"`
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data"`
	Code     int                    `json:"code"`
}

type Response struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func (r *Response) JSON(statusCode int, data *BizData) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.StatusCode = statusCode
	r.Body = string(marshal)
	return nil
}

func (r *Response) OKWithJSON(data interface{}, meta ...map[string]interface{}) error {
	bz := &BizData{Msg: "ok", Code: 0, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.JSON(200, bz)
}
