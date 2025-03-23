package runner

import "fmt"

func (r *Runner) routerInfo(ctx *Context) {
	//apphub _get_router_info /array/diff GET
	router := r.args[2]
	method := r.args[3]
	worker, exist := r.getRouterWorker(router, method)
	if !exist {
		ctx.Response.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s not exist", method, router))
		return
	}
	params, err := worker.Config.GetParams()
	if err != nil {
		ctx.Response.FailWithJSON(nil, fmt.Sprintf("method:%s router:%s GetParams err:%s", method, router, err))
		return
	}
	ctx.Response.OKWithJSON(params)
}
