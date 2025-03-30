package runner

var r = New()

func Post(router string, fn func(ctx *HttpContext), config ...*Config) {
	if router[0] != '/' {
		router = "/" + router
	}
	r.Post(router, fn, config...)
}
func Get(router string, fn func(ctx *HttpContext), config ...*Config) {
	if router[0] != '/' {
		router = "/" + router
	}
	r.Get(router, fn, config...)
}

func Run() error {
	return r.Run()
}
