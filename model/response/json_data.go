package response

import "github.com/bytedance/sonic"

type jsonData struct {
	response  *Response
	TraceID   string                 `json:"trace_id"`
	MetaData  map[string]interface{} `json:"meta_data"` //sdk 层
	Code      int                    `json:"code"`
	Msg       string                 `json:"msg"`
	Data      interface{}            `json:"data"`
	BuildData string                 `json:"build_data"`
}

func (j *jsonData) DataType() DataType {
	return DataTypeJSON
}

func (j *jsonData) Build() error {
	marshal, err := sonic.Marshal(j)
	if err != nil {
		return err
	}
	j.BuildData = string(marshal)
	//j.response.data = append(j.response.data, j)
	return nil
}
