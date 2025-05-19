package response

import "github.com/pkg/errors"

const (
	RenderTypeForm    = "form"
	RenderTypeJSON    = "json"
	RenderTypeTable   = "table"
	RenderTypeFiles   = "files"
	RenderTypeEcharts = "echarts"
)

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
