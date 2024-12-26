package runner

type Context struct {
	runner       *Runner
	HttpRequest  *Request
	HttpResponse *Response
}
