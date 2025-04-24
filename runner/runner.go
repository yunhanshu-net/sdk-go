package runner

import (
	"context"
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
	"runtime/debug"
	"strings"
	"time"
)

func New() *Runner {
	return &Runner{
		idle:         5,
		lastHandelTs: time.Now(),
		//handelFunctions: make(map[string]*Worker),
		routerMap: make(map[string]*routerInfo),
		down:      make(chan struct{}, 1),
	}
}

type Runner struct {
	isDebug       bool
	detail        *model.Runner
	uuid          string
	args          []string
	idle          int64
	lastHandelTs  time.Time
	isClosed      bool
	natsConn      *nats.Conn
	natsSubscribe *nats.Subscription
	routerMap     map[string]*routerInfo
	down          chan struct{}
}

func (r *Runner) init(args []string) error {
	r.args = args
	r.detail = &model.Runner{}
	split := strings.Split(r.args[0], "_")
	if len(split) > 1 {
		r.detail.User = strings.ReplaceAll(split[0], "./", "")
		r.detail.Name = split[1]
	}
	fmt.Println("detail:", r.detail)

	runtime.GOMAXPROCS(2)
	r.get("/_env", env)
	r.get("/_ping", ping)
	//r.get("/_router_info", r.routerInfo)
	//r.get("/_router_list_info", r.routerListInfo)
	r.get("/_get_api_infos", r.getApiInfos)
	r.get("/_get_api_info", r.getApiInfo) // 添加新的路由
	r.post("/_callback", r.callback)

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
		go func() {
			err = r.connectNats()
			if err != nil {
				logrus.Infof("connect err:%s", err.Error())
				panic(err)
			}
		}()
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

func (r *Runner) getRouter(router string, method string) (worker *routerInfo, exist bool) {
	worker, ok := r.routerMap[fmtKey(router, method)]
	if ok {
		return worker, true
	}
	return nil, false
}

func (r *Runner) runRequest(ctx0 context.Context, req *request.Request) (*response.Data, error) {
	//worker, exist := r.getRouterWorker(ctx.Request.Route, ctx.Request.Method)
	router, exist := r.getRouter(req.Route, req.Method)
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			stack := debug.Stack()
			// 增加更详细的错误信息输出
			fmt.Printf("具体错误: %v\n", errPanic)
			logrus.Errorf("runRequest panic err:%s req:%+v stack:%s", errPanic, req, string(stack))
			fmt.Println(string(stack))
		}
	}()

	if !exist {
		marshal, _ := json.Marshal(r.routerMap)
		fmt.Printf("handels: %s\n", string(marshal))
		fmt.Printf("method:%s router:%s not found\n", req.Method, req.Route)
		return nil, fmt.Errorf("method:%s router:%s not found\n", req.Method, req.Route)
	}

	//marshal, err := sonic.Marshal(req.Body)
	//if err != nil {
	//	return nil, err
	//}
	now := time.Now()
	_, rsp, err := router.call(ctx0, req.Body)
	if err != nil {
		logrus.Errorf("runRequest err:%s", err.Error())
		return nil, err
	}
	since := time.Since(now)
	if rsp.MetaData == nil {
		rsp.MetaData = make(map[string]interface{})
	}
	rsp.MetaData["cost"] = since.String()
	//ctx.Response = rsp
	//todo 判断是否需要reset body

	return rsp, nil
}

func (r *Runner) run(req *request.RunnerRequest) {

	ctx := context.Background()
	resp, err := r.runRequest(ctx, req.Request)
	if err != nil {
		panic(err)
	}
	marshal, err := sonic.Marshal(resp)
	if err != nil {
		panic(err)
	}
	//todo 这里只是用来测试
	//jsonx.SaveFile("./out.json", httpContext.Response)
	fmt.Println("<Response>" + string(marshal) + "</Response>")
}
