package request

type NoData struct{}

type RouterInfo struct {
	Router string `json:"router"`
	Method string `json:"method"`
}

type ApiInfoRequest struct {
	Router string `json:"router"` // API路由路径
	Method string `json:"method"` // HTTP方法（GET/POST）
}
