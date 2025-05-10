package response

import (
	"encoding/json"
	"net/http"
)

const (
	successMsg  = "ok"
	successCode = 0
)

type Form interface {
	Builder
}

func (r *Data) Form(data interface{}) Form {
	r.RenderType = RenderTypeForm
	bz := &formData{TraceID: r.TraceID, Msg: successMsg, Code: successCode, Data: data, response: r}
	r.StatusCode = http.StatusOK
	return bz
}

func (r *Data) FailWithJSON(data interface{}, msg string, meta ...map[string]interface{}) error {
	r.RenderType = RenderTypeForm

	bz := &formData{TraceID: r.TraceID, Msg: msg, Code: -1, Data: data, DataType: RenderTypeForm, response: r}
	if len(meta) > 0 {
		bz.MetaData = meta[0]
	}
	return bz.Build()
}

type formResp struct {
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	DataType RenderType             `json:"data_type"`
	TraceID  string                 `json:"trace_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Data     interface{}            `json:"data"`
}

type formData struct {
	response  *Data
	DataType  RenderType             `json:"data_type"`
	TraceID   string                 `json:"trace_id"`
	MetaData  map[string]interface{} `json:"meta_data"` //sdk 层
	Code      int                    `json:"code"`
	Msg       string                 `json:"msg"`
	Data      interface{}            `json:"data"`
	buildData string
}

func (j *formData) GetRenderType() RenderType {
	return RenderTypeForm
}

func (j *formData) Build() error {
	j.response.RenderType = j.GetRenderType()
	j.response.Multiple = false
	r := &formResp{
		Code:     j.Code,
		Msg:      j.Msg,
		Data:     j.Data,
		DataType: RenderTypeForm,
		TraceID:  j.TraceID,
		MetaData: j.MetaData,
	}
	marshal, err := json.Marshal(r)
	if err != nil {
		return err
	}
	j.response.Body = string(marshal)
	return nil
}

func (j *formData) BuildJSON() string {
	j.response.RenderType = j.GetRenderType()
	j.response.Multiple = false
	j.response.Body = j.buildData
	return j.buildData
}
