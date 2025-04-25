package response

import (
	"github.com/bytedance/sonic"
	"net/http"
)

const (
	successMsg  = "ok"
	successCode = 0
)

type JSON interface {
	Builder
}

func (r *Data) JSON(data interface{}) JSON {
	r.DataType = DataTypeJSON
	bz := &jsonData{Msg: successMsg, Code: successCode, Data: data, response: r}
	r.StatusCode = http.StatusOK
	return bz
}

func (r *Data) FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error {
	r.DataType = DataTypeJSON

	bz := &jsonData{Msg: msg, Code: -1, Data: data, DataType: DataTypeJSON, response: r}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return bz.Build()
}

type formResp struct {
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	DataType DataType               `json:"data_type"`
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Data     interface{}            `json:"data"`
}

type jsonData struct {
	response  *Data
	DataType  DataType               `json:"data_type"`
	TraceID   string                 `json:"trace_id"`
	MetaData  map[string]interface{} `json:"meta_data"` //sdk 层
	Code      int                    `json:"code"`
	Msg       string                 `json:"msg"`
	Data      interface{}            `json:"data"`
	buildData string
}

func (j *jsonData) GetDataType() DataType {
	return DataTypeJSON
}

func (j *jsonData) Build() error {
	j.response.DataType = j.GetDataType()
	j.response.Multiple = false
	r := &formResp{
		Code:     j.Code,
		Msg:      j.Msg,
		Data:     j.Data,
		DataType: DataTypeJSON,
		TraceID:  j.TraceID,
		MetaData: j.MetaData,
	}
	marshal, err := sonic.Marshal(r)
	if err != nil {
		return err
	}
	j.response.Body = string(marshal)
	return nil
}

func (j *jsonData) BuildJSON() string {
	j.response.DataType = j.GetDataType()
	j.response.Multiple = false
	j.response.Body = j.buildData
	return j.buildData
}
