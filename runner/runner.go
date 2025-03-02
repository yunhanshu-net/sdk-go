package runner

import (
	"context"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"sync"
)

var (
	defaultMaxIdealTime int64 = 5 //秒
)

type TransportConfig struct {
	TransportType string            `json:"transport_type"`
	Metadata      map[string]string `json:"metadata"`
	WorkPath      string            `json:"work_path"`
	RunnerType    string            `json:"runner_type"`
	Version       string            `json:"version"`
	Route         string            `json:"route"`  //命令
	User          string            `json:"user"`   //软件所属的用户
	Runner        string            `json:"runner"` //软件名
	OssPath       string            `json:"oss_path"`
	StartArgs     []string          `json:"start_args"`
}
type Worker struct {
	Handel []func(ctx *Context)
	Path   string
	Method string
	Config *Config
}

type Runner struct {
	Transport       Transport
	idle            int
	isDebug         bool
	transportConfig *TransportConfig
	handelFunctions map[string]*Worker
	notFound        func(ctx *Context)
	wg              *sync.WaitGroup
	lastHandelTs    int64
	isKeepAlive     bool

	conn Conn
	args []string

	//nats            *nats.Conn
	//contextChan chan *Context
	//sub             *nats.Subscription

	requestCh chan *request.Request
	//exit            <-chan struct{}
	exitCtx context.Context
	exit    context.CancelFunc
}

func (r *Runner) Close() {
	r.conn.Close()
}

func (r *Runner) runRequest(method string, router string, ctx *Context) error {
	worker := r.handelFunctions[router+"."+method]
	for _, fn := range worker.Handel {
		fn(ctx)
	}
	return nil
}
