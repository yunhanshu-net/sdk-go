package request

type NoData struct{}

type RouterInfo struct {
	Router string `json:"router"`
	Method string `json:"method"`
}
