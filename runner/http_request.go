package runner

func (r *Runner) Post(router string, handelFunc func(ctx *HttpContext), config ...*Config) {
	_, ok := r.handelFunctions[r.fmtHandelKey(router, "POST")]
	if !ok {
		worker := &Worker{
			Handel: []func(ctx *HttpContext){handelFunc},
			Method: "POST",
			Path:   router,
			Config: &Config{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}
		r.handelFunctions[r.fmtHandelKey(router, "POST")] = worker
	} else {
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}

}
func (r *Runner) Get(router string, handelFunc func(ctx *HttpContext), config ...*Config) {
	_, ok := r.handelFunctions[r.fmtHandelKey(router, "GET")]
	if !ok {
		worker := &Worker{
			Handel: []func(ctx *HttpContext){handelFunc},
			Method: "GET",
			Path:   router,
			Config: &Config{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}

		r.handelFunctions[r.fmtHandelKey(router, "GET")] = worker
	} else {
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}
}
