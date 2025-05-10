package response

type RenderType string

const (
	RenderTypeForm    = "form"
	RenderTypeJSON    = "json"
	RenderTypeTable   = "table"
	RenderTypeFiles   = "files"
	RenderTypeEcharts = "echarts"
)

type Builder interface {
	Build() error
}

type anyData interface {
	Builder
	BuildJSON() string
	GetRenderType() RenderType
}
