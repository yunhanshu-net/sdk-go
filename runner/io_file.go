package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type fileIo struct {
	filePath string
}

func (f *fileIo) Input(runner *Runner) *request.Request {

	//command := runner.args[1]
	//jsonFileName := runner.args[2]

	//TODO implement me
	panic("implement me")
}

func (f *fileIo) Output(runner *Runner, rsp *response.Response) {
	//TODO implement me
	panic("implement me")
}
