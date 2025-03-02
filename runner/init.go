package runner

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"runtime"
	"time"
)

func (r *Runner) connect(req *TransportConfig) {
	r.transportConfig = req
	//todo 长连接
	if req.TransportType == "nats" || req.TransportType == "" {
		trs, err := newTransportNats(req)
		if err != nil {
			fmt.Println("newTransportNats err", err.Error())
			panic(err)
		}
		fmt.Printf("连接nats成功\n")
		r.lastHandelTs = time.Now().Unix()
		r.Transport = trs
	}
}

func (r *Runner) run(req *Request) {
	ctx := &Context{
		transportConfig: req.TransportConfig,
		Request:         req.Request,
		Response:        &response.Response{},
	}
	err := r.runRequest(req.Request.Method, req.Request.Route, ctx)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(ctx.Response)
	if err != nil {
		panic(err)
	}
	fmt.Println("<Response>" + string(marshal) + "</Response>")
}

func (r *Runner) init(args []string) {
	r.args = args
	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)
	req, err := r.getRequest(r.args[2])
	if err != nil {
		panic(err)
	}
	if r.args[1] == "_connect" { //长连接
		//todo 长连接
		r.connect(req.TransportConfig)
	} else { //说明是单次执行
		r.run(req)
		return
	}
	r.listen()
}

func (r *Runner) GetCommand() string {
	return r.args[1]
}

func (r *Runner) getRequest(filePath string) (*Request, error) {
	var req Request
	req.Request = new(request.Request)
	err := jsonx.UnmarshalFromFile(filePath, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Runner) listen() {
	//timeout := 2 * time.Second
	//idleTimer := time.NewTimer(timeout)
	//defer idleTimer.Stop()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	//maxIdealTime := time.NewTicker(time.Second)
	for {
		// 每次处理消息后重置定时器
		//idleTimer.Reset(timeout)
		select {
		case msg, ok := <-r.Transport.ReadMessage():
			if !ok {
				fmt.Printf("runner:%s ReadMessage chan closed \n", r.info())
				return
			}
			r.lastHandelTs = time.Now().Unix()
			go r.handelMsg(msg, r)

		//case <-idleTimer.C:
		case <-ticker.C:
			if r.idle > 0 {
				ts := time.Now().Unix()
				if (ts - r.lastHandelTs) > int64(r.idle) { //超过指定空闲时间的话需要释放进程
					err := r.Transport.Close()
					if err != nil {
						panic(err)
					}
					return
				}
			}

		default:

		}
	}
}

func (r *Runner) handelMsg(transportMsg *TransportMsg, runner *Runner) {
	method := transportMsg.Headers.Get("method")
	route := transportMsg.Headers.Get("route")
	var reqMsg request.Request
	err1 := json.Unmarshal(transportMsg.Data, &reqMsg)
	if err1 != nil {
		panic(err1)
	}
	if method == "" {
		method = reqMsg.Method
	}
	if route == "" {
		route = reqMsg.Route
	}
	ctx := &Context{Request: &reqMsg, Response: &response.Response{}, transportConfig: r.transportConfig}
	err := runner.runRequest(method, route, ctx)
	if err1 != nil {
		panic(err)
	}
	t := newTransportMsg(transportMsg)
	marshal, err1 := json.Marshal(ctx.Response)
	if err1 != nil {
		panic(err)
	}
	t.Data = marshal
	err1 = transportMsg.Reply(t)
	if err1 != nil {
		panic(err1)
	}
}

func (r *Runner) info() string {
	return fmt.Sprintf("%s.%s.%s", r.transportConfig.User, r.transportConfig.Runner, r.transportConfig.Version)
}
