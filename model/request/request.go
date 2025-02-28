package request

type Request struct {
	TraceID string                 `json:"trace_id"`
	Route   string                 `json:"route"`
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"` //请求json
	FileMap map[string][]string    `json:"file_map"`
}
