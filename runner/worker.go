package runner

import "strings"

type Worker struct {
	Handel []func(ctx *HttpContext) error `json:"-"`
	Path   string
	Method string
	Config *ApiConfig
}

// IsDefaultRouter _开头的路由是默认路由
func (w *Worker) IsDefaultRouter() bool {
	return strings.HasPrefix(strings.TrimPrefix(w.Path, "/"), "_")
}
