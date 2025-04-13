package runner

func (r *Runner) get(router string, handel interface{}, config ...*ApiConfig) {
	key := fmtKey(router, "GET")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:    key,
			Handel: handel,
			Method: "GET",
			Router: router,
			Config: &ApiConfig{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) post(router string, handel interface{}, config ...*ApiConfig) {
	key := fmtKey(router, "POST")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:    key,
			Handel: handel,
			Method: "POST",
			Router: router,
			Config: &ApiConfig{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}
