package runner

import (
	"github.com/nats-io/nats.go"
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
	nats        *nats.Conn
	sub         *nats.Subscription
	isKeepAlive bool
	contextChan chan *Context
	wg          *sync.WaitGroup
	exit        <-chan struct{}
	info        Info
}

func (r *Runner) Close() {
	r.sub.Unsubscribe()
	r.nats.Close()
}
