package request

import "github.com/yunhanshu-net/pkg/dto/runnerproject"

type RunFunctionReq struct {
	RunnerID string                `json:"runner_id"`
	Runner   *runnerproject.Runner `json:"runner"`
	TraceID  string                `json:"trace_id"`
	Router   string                `json:"router"`
	Method   string                `json:"method"`
	Headers  map[string]string     `json:"headers"`
	BodyType string                `json:"body_type"`
	Body     interface{}           `json:"body"`
	UrlQuery string                `json:"url_query"`
}
