package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

func env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"ping": "pong"}).Build()
}

func (r *Runner) routerListInfo(ctx *Context, req *request.NoData, resp response.Response) error {
	functions := r.routerMap
	var configs []*ApiConfig
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		worker.Config.Method = worker.Method
		worker.Config.Router = worker.Router
		if worker.Config != nil {
			if worker.Config.Request != nil {
				params, err := worker.Config.getParams(worker.Config.Request, "in")
				if err != nil {
					continue
				}
				worker.Config.ParamsIn = params
			}

			if worker.Config.Response != nil {
				params, err := worker.Config.getParams(worker.Config.Response, "out")
				if err != nil {
					continue
				}
				worker.Config.ParamsOut = params
			}
			configs = append(configs, worker.Config)
		}
	}

	return resp.JSON(configs).Build()
}
