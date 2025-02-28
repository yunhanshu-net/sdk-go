package runner

type Request struct {
	Url      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"` //请求json
	FileList []string          `json:"file_list"`
	Info     Info              `json:"info"`
}
