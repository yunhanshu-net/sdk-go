package request

type Request struct {
	TraceID string                 `json:"trace_id"`
	Route   string                 `json:"route"`
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"` //请求json
	FileMap map[string][]string    `json:"file_map"`
}

type DebugRequest struct {
	Request

	User    string `json:"user"`
	Runner  string `json:"runner"`
	Version string `json:"version"`
}

type Runner struct {
	Name    string `json:"name"`
	User    string `json:"user"`
	Version string `json:"version"`
}
type RunnerRequest struct {
	Runner          *Runner                `json:"runner"`
	TransportConfig *TransportConfig       `json:"transport_config"`
	Metadata        map[string]interface{} `json:"metadata"`
	Request         *Request               `json:"request"`
}

type TransportConfig struct {
	Type string `json:"type"`
}
