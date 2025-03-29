package runner

type Worker struct {
	Handel []func(ctx *HttpContext) `json:"-"`
	Path   string
	Method string
	Config *Config
}
