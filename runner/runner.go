package runner

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	v2 "github.com/yunhanshu-net/sdk-go/model/response/v2"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"runtime"
	"strings"
	"sync"
	"time"
)

func New() *Runner {
	return &Runner{
		idle:            5,
		lastHandelTs:    time.Now(),
		handelFunctions: make(map[string]*Worker),
		wg:              &sync.WaitGroup{},
	}
}

type Runner struct {
	isDebug bool
	detail  *model.Runner
	uuid    string
	rpcSrv  *server.Server
	conn    *nats.Conn
	sub     *nats.Subscription
	wg      *sync.WaitGroup
	args    []string

	idle            int64
	lastHandelTs    time.Time
	handelFunctions map[string]*Worker

	down <-chan struct{}
}

func (r *Runner) fmtHandelKey(router string, method string) string {
	if !strings.HasPrefix(router, "/") {
		router = "/" + router
	}
	router = strings.TrimSuffix(router, "/")
	return router + "." + strings.ToUpper(method)

}

func (r *Runner) init(args []string) error {
	r.args = args
	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)
	r.Get("/_router_info", r.routerInfo)
	r.Get("/_router_list_info", r.routerListInfo)
	var err error
	var req = new(request.RunnerRequest)
	req, err = r.getRequest(r.args[2])
	if err != nil {
		panic(err)
	}

	if req != nil {
		r.detail = req.Runner
		if r.uuid == "" {
			r.uuid = req.UUID
		}
	}

	if r.args[1] == "_connect" { //长连接
		if req.TransportConfig != nil && req.TransportConfig.IdleTime != 0 { //最大空闲时间
			r.idle = int64(req.TransportConfig.IdleTime)
		}
		//todo 长连接
		err = r.connectRpc()
		if err != nil {
			logrus.Infof("connect err:%s", err.Error())
			panic(err)
		}
		r.listen()
		logrus.Infof("uuid:%s listen stop\n", r.uuid)
		return nil
	} else { //说明是单次执行
		r.run(req)
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

func (r *Runner) getRouterWorker(router string, method string) (worker *Worker, exist bool) {
	if router[0] != '/' {
		router = "/" + router
	}
	worker, ok := r.handelFunctions[r.fmtHandelKey(router, method)]
	if ok {
		return worker, true
	}
	return nil, false
}
func (r *Runner) runRequest(ctx *HttpContext) error {
	worker, exist := r.getRouterWorker(ctx.Request.Route, ctx.Request.Method)
	if !exist {
		marshal, _ := json.Marshal(r.handelFunctions)
		fmt.Printf("handels: %s\n", string(marshal))
		fmt.Printf("method:%s router:%s not found\n", ctx.Request.Method, ctx.Request.Route)
		return fmt.Errorf("method:%s router:%s not found\n", ctx.Request.Method, ctx.Request.Route)
	}
	for _, fn := range worker.Handel {
		fn(ctx)
	}
	return nil
}

func (r *Runner) run(req *request.RunnerRequest) {
	//ctx := &Context{
	//	req:      req,
	//	runner:   req.Runner,
	//	Request:  req.Request,
	//	ResponseData: &response.ResponseData{},
	//}

	rsp := &v2.ResponseData{
		MetaData: make(map[string]interface{}),
	}
	now := time.Now()
	httpContext := &HttpContext{
		Request:  req.Request,
		runner:   req.Runner,
		Response: rsp,
	}

	err := r.runRequest(httpContext)
	if err != nil {
		panic(err)
	}
	rsp.MetaData["func_cost"] = time.Since(now).String()
	marshal, err := json.Marshal(httpContext.Response)
	if err != nil {
		panic(err)
	}
	//todo 这里只是用来测试
	jsonx.SaveFile("./out.json", httpContext.Response)
	fmt.Println("<Response>" + string(marshal) + "</Response>")
}
