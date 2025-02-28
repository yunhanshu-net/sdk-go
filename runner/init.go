package runner

import (
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"os"
	"runtime"
)

func (r *Runner) init() {
	r.args = os.Args

	//判断是冷启动还是建立长链接
	//如果是运行指令的话

	request, err := r.getRequest()
	if err != nil {
		panic(err)
	}
	r.info = &request.Info
	if r.args[1] == "_connect" {
		//todo 长连接

	}
	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)

}

func (r *Runner) GetCommand() string {
	return r.args[1]
}

func (r *Runner) getRequest() (*Request, error) {
	var req Request
	err := jsonx.UnmarshalFromFile(r.args[2], &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
