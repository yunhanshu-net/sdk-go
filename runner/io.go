package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type IO interface {
	Input(runner *Runner) (*request.Request, error)
	Output(runner *Runner, rsp *response.Response) error
}
