package runner

var configMap = map[string]string{}

func GetConfig[T any](ctx *Context) T {
	var t T
	return t
}

// InitConfig todo 这里加载指定接口的配置文件
func InitConfig[T any](ctx *Context, conf T) T {
	var t T
	return t
}
