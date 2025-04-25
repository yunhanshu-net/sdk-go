package request

type NoData struct{}

type ApiInfoRequest struct {
	Router string `json:"router" form:"router"` // API路由路径
	Method string `json:"method" form:"method"` // HTTP方法（GET/POST）
}
