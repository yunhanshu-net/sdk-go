package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

func (r *Runner) routerInfo(ctx *Context, req *routerInfo, resp response.Response) error {

	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		return resp.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s not exist", req.Method, req.Router))
	}
	params, err := worker.Config.GetParams()
	if err != nil {
		return resp.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s GetParams err:%s", req.Method, req.Router, err))
	}
	return resp.JSON(params).Build()
}
