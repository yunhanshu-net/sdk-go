package runner

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"os"
	"runtime"
)

func (r *Runner) init() {
	r.args = os.Args

	//判断是冷启动还是建立长链接
	//如果是运行指令的话

	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)

	req, err := r.getRequest()
	if err != nil {
		panic(err)
	}
	if r.args[1] == "_connect" {
		//todo 长连接
		if req.TransportConfig.TransportType == "nats" || req.TransportConfig.TransportType == "" {
			trs, err := newTransportNats(req.TransportConfig)
			if err != nil {
				panic(err)
			}
			r.Transport = trs
		}
	} else { //说明是单次执行
		ctx := &Context{
			transportConfig: req.TransportConfig,
			Request:         req.Request,
			Response:        &response.Response{},
		}
		err := r.RunRequest(req.Request.Method, req.Request.Route, ctx)
		if err != nil {
			panic(err)
		}
		marshal, err := json.Marshal(ctx.Response)
		if err != nil {
			panic(err)
		}
		fmt.Println("<Response>" + string(marshal) + "</Response>")
		return
	}

	for {
		select {
		case msg := <-r.Transport.ReadMessage():
			go func(transportMsg *TransportMsg, runner *Runner) {
				method := transportMsg.Headers.Get("method")
				route := transportMsg.Headers.Get("route")
				var reqMsg request.Request
				err1 := json.Unmarshal(transportMsg.Data, &reqMsg)
				if err1 != nil {
					panic(err1)
				}
				ctx := &Context{Request: &reqMsg, Response: &response.Response{}, transportConfig: r.info}
				err = runner.RunRequest(method, route, ctx)
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

			}(msg, r)
		case <-r.Transport.Wait():
			return

		default:

		}
	}

}

func (r *Runner) GetCommand() string {
	return r.args[1]
}

func (r *Runner) getRequest() (*Request, error) {
	var req Request
	req.Request = new(request.Request)
	err := jsonx.UnmarshalFromFile(r.args[2], &req.Request)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
