package runner

func env(ctx *HttpContext) error {
	return ctx.Response.JSON(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func ping(ctx *HttpContext) error {
	return ctx.Response.JSON(map[string]string{"ping": "pong"}).Build()
}

func (r *Runner) routerListInfo(ctx *HttpContext) error {
	functions := r.handelFunctions
	var configs []*ApiConfig
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		worker.Config.Method = worker.Method
		worker.Config.Router = worker.Path
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

	return ctx.Response.JSON(configs).Build()
}
