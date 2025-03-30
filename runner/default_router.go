package runner

import "fmt"

type routerInfo struct {
	Router string `json:"router"`
	Method string `json:"method"`
}

func (r *Runner) routerInfo(ctx *HttpContext) {

	var req routerInfo
	err := ctx.Request.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}
	//apphub _get_router_info /array/diff GET
	//router := r.args[2]
	//method := r.args[3]
	worker, exist := r.getRouterWorker(req.Router, req.Method)
	if !exist {
		ctx.Response.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s not exist", req.Method, req.Router))
		return
	}
	params, err := worker.Config.GetParams()
	if err != nil {
		ctx.Response.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s GetParams err:%s", req.Method, req.Router, err))
		return
	}
	ctx.Response.JSON(params).Build()
}
