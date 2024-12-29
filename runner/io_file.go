package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"os"
)

func newFileIo(runner *Runner) *fileIo {
	return &fileIo{
		filePath: runner.args[2],
	}
}

type fileIo struct {
	filePath string
}

func (f *fileIo) Input(runner *Runner) (*request.Request, error) {
	file, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, err
	}

	//command := runner.args[1]
	//jsonFileName := runner.args[2]

	//TODO implement me
	panic("implement me")
}

func (f *fileIo) Output(runner *Runner, rsp *response.Response) error {
	//TODO implement me
	panic("implement me")
}
