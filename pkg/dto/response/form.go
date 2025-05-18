package response

import "errors"

const (
	RenderTypeForm    = "form"
	RenderTypeJSON    = "json"
	RenderTypeTable   = "table"
	RenderTypeFiles   = "files"
	RenderTypeEcharts = "echarts"
)

type formData struct {
	resp *RunFunctionResp
	Data interface{} `json:"data"`
}

func newForm(data interface{}, resp *RunFunctionResp) Form {
	return &formData{Data: data, resp: resp}
}

func build(resp *RunFunctionResp, data interface{}, renderType string) error {
	if resp == nil {
		return errors.New("resp is nil")
	}
	if resp.Data != nil {
		resp.Multiple = true
		resp.DataList = append(resp.DataList, resp.Data)
		resp.DataList = append(resp.DataList, data)
		resp.DataType = resp.DataType + "," + renderType
	}
	resp.DataType = renderType
	resp.Data = data
	return nil
}
func (d *formData) Build() error {
	return build(d.resp, d.Data, RenderTypeForm)
}
