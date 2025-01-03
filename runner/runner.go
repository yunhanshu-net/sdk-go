package runner

import (
	"context"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"sync"
)

type Info struct {
	WorkPath   string `json:"work_path"`
	RunnerType string `json:"runner_type"`
	Version    string `json:"version"`
	Command    string `json:"command"` //命令
	User       string `json:"user"`    //软件所属的用户
	Soft       string `json:"soft"`    //软件名
	OssPath    string `json:"oss_path"`
}

type Runner struct {
	lastHandelTs int64
	conn         Conn
	args         []string
	isKeepAlive  bool
	info         *Info
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
