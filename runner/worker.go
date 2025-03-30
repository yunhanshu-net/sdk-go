package runner

import "strings"

type Worker struct {
	Handel []func(ctx *HttpContext) `json:"-"`
	Path   string
	Method string
	Config *Config
}

// IsDefaultRouter _开头的路由是默认路由
func (w *Worker) IsDefaultRouter() bool {
	return strings.HasPrefix(strings.TrimPrefix(w.Path, "/"), "_")
}
