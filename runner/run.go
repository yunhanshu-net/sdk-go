package runner

import (
	"context"
	"sync"
)

const (
	commandConnect = "_connect_"
)

func (r *Runner) listenAndServe() error {

	//jsonFileName := os.Args[2]

	command := r.GetCommand()
	conn := NewConn(r)
	err := conn.Connect()
	if err != nil {
		return err
	}
	r.conn = conn
	go r.readLoop()
	switch command {
	case commandConnect: //建立长连接
		//r.isKeepAlive=true

	default:
		defer r.exit()
	}

	select {
	case <-r.exitCtx.Done():
		//todo 释放资源
		if r.isKeepAlive {
			r.wg.Wait()
		}
		return nil
	}

}

func (r *Runner) Run() {
	r.init()
	command := r.GetCommand()
	if command == commandConnect {
		r.isKeepAlive = true
	}
	r.listenAndServe()
	//单核心执行，防止程序把资源占用尽
}

func New() *Runner {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Runner{
		exitCtx:         ctx,
		handelFunctions: make(map[string]*Worker),
		exit:            cancelFunc,
		wg:              &sync.WaitGroup{},
	}
}
