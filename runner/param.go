package runner

type RequestDesc struct {
	Desc string `json:"desc"` //请求参数的描述信息 例如：
}

type Params interface {
	Desc() *RequestDesc
}
