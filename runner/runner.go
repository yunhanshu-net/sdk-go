package runner

import (
	"context"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"sync"
)

type ContextInfo struct {
	Metadata   map[string]string `json:"metadata"`
	WorkPath   string            `json:"work_path"`
	RunnerType string            `json:"runner_type"`
	Version    string            `json:"version"`
	Command    string            `json:"command"` //命令
	User       string            `json:"user"`    //软件所属的用户
	Soft       string            `json:"soft"`    //软件名
	OssPath    string            `json:"oss_path"`
	StartArgs  []string          `json:"start_args"`
}

type Runner struct {
	Transport Transport

	lastHandelTs int64
	conn         Conn
	args         []string
	isKeepAlive  bool
	info         *ContextInfo
	//nats            *nats.Conn
	//contextChan chan *Context
	//sub             *nats.Subscription
	handelFunctions map[string]*Worker
	wg              *sync.WaitGroup
	requestCh       chan *request.Request
	//exit            <-chan struct{}
	exitCtx  context.Context
	exit     context.CancelFunc
	notFound func(ctx *Context)
}

func (r *Runner) Close() {
	r.conn.Close()
}
