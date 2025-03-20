package runner

import (
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"runtime"
	"sync"
	"time"
)

func New() *Runner {
	fmt.Printf("new --------")
	return &Runner{
		idle:            5,
		lastHandelTs:    time.Now(),
		handelFunctions: make(map[string]*Worker),
		wg:              &sync.WaitGroup{},
	}
}

type Runner struct {
	detail *model.Runner
	uuid   string
	conn   *nats.Conn
	sub    *nats.Subscription
	wg     *sync.WaitGroup
	args   []string

	idle            int64
	lastHandelTs    time.Time
	handelFunctions map[string]*Worker
}

func (r *Runner) init(args []string) error {
	r.args = args
	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)
	req, err := r.getRequest(r.args[2])
	if err != nil {
		panic(err)
	}

	r.detail = req.Runner
	if r.uuid == "" {
		r.uuid = req.UUID
	}

	if r.args[1] == "_connect" { //长连接
		//todo 长连接
		err = r.connect()
		if err != nil {
			logrus.Infof("connect err:%s", err.Error())
			panic(err)
		}
		r.listen()
		logrus.Infof("uuid:%s listen stop\n", r.uuid)
		return nil
	} else { //说明是单次执行
		r.run(req)
		logrus.Infof("uuid:%s run stop\n", r.uuid)
		return nil
	}

	return nil
}

func (r *Runner) getRequest(filePath string) (*request.RunnerRequest, error) {
	var req request.RunnerRequest
	req.Request = new(request.Request)
	err := jsonx.UnmarshalFromFile(filePath, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Runner) runRequest(method string, router string, ctx *Context) error {
	if router[0] != '/' {
		router = "/" + router
	}

	worker, ok := r.handelFunctions[router+"."+method]
	if !ok {
		marshal, _ := json.Marshal(r.handelFunctions)
		fmt.Printf("handels: %s\n", string(marshal))
		fmt.Printf("method:%s router:%s not found\n", method, router)
		return fmt.Errorf("method:%s router:%s not found\n", method, router)
	}
	for _, fn := range worker.Handel {
		fn(ctx)
	}
	return nil
}

func (r *Runner) run(req *request.RunnerRequest) {
	ctx := &Context{
		req:      req,
		runner:   req.Runner,
		Request:  req.Request,
		Response: &response.Response{},
	}
	err := r.runRequest(req.Request.Method, req.Request.Route, ctx)
	if err != nil {
		panic(err)
	}
	marshal, err := sonic.Marshal(ctx.Response)
	if err != nil {
		panic(err)
	}
	fmt.Println("<Response>" + string(marshal) + "</Response>")
}
