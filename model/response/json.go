package response

//type IData interface {
//	DataType() DataType
//	Build() error
//}

//type JSON interface {
//
//
//}

type JSON struct {
	response *Response
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"` //sdk 层
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	Data     interface{}            `json:"data"`
}

func (j *JSON) DataType() DataType {
	return DataTypeJSON
}

func (r *Response) JSON(data interface{}) *JSON {
	r.DataType = DataTypeJSON
	bz := &JSON{Msg: "ok", Code: 0, Data: data, response: r}

	r.StatusCode = 200
	return bz
}

func (j *JSON) Build() error {
	j.response.DataType = DataTypeJSON
	//j.response.data = append(j.response.data, &jsonData{
	//	TraceID: j.TraceID,
	//	Data:    j.Data,
	//	Code:    j.Code,
	//	Msg:     j.Msg,
	//})
	return nil
}

func (j *JSON) Fail(msg string) *JSON {
	j.Code = -1
	j.Msg = msg
	return j
}

func (r *Response) json(statusCode int, data *JSON) error {
	r.StatusCode = statusCode
	r.Body = data
	return nil
}

func (r *Response) OKWithJSON(data interface{}, meta ...map[string]interface{}) error {
	r.DataType = DataTypeJSON
	bz := &JSON{Msg: "ok", Code: 0, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
func (r *Response) FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error {
	r.DataType = DataTypeJSON
	bz := &JSON{Msg: msg, Code: -1, Data: data}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return r.json(200, bz)
}
