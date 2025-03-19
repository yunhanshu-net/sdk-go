package runner

var r = New()

func Post(router string, fn func(ctx *Context), config ...*Config) {
	r.Post(router, fn)
}
func Get(router string, fn func(ctx *Context), config ...*Config) {
	r.Post(router, fn)
}

func Run() error {
	return r.Run()
}
