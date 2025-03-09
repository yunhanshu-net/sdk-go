package v2

var r = New()

func Post(router string, fn func(ctx *Context)) {
	r.Post(router, fn)
}
func Get(router string, fn func(ctx *Context)) {
	r.Post(router, fn)
}

func Run() error {
	return r.Run()
}
