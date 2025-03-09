package v2

type Worker struct {
	Handel []func(ctx *Context)
	Path   string
	Method string
	Config *Config
}
