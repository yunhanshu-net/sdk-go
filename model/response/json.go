package response

type IDataType interface {
	DataType() DataType
}

type JSON struct {
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"` //sdk 层
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data"`
}

func (j *JSON) DataType() DataType {
	return DataTypeJSON
}

func (r *Response) json(statusCode int, data *JSON) error {
	r.StatusCode = statusCode
	r.Body = data
	return nil
}

func (r *Response) OKWithJSON(data interface{}, meta ...map[string]interface{}) error {
	bz := &JSON{Msg: "ok", Code: 0, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
func (r *Response) FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error {
	bz := &JSON{Msg: msg, Code: -1, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
