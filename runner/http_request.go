package runner

func (r *Runner) Post(router string, handelFunc func(ctx *Context), opts ...Option) {
	_, ok := r.handelFunctions[router]
	if !ok {
		worker := &Worker{
			Handel: []func(*Context){handelFunc},
			Method: "POST",
			Path:   router,
			Config: &Config{},
		}
		if len(opts) > 0 {
			for _, opt := range opts {
				opt(worker.Config)
			}
		}
		r.handelFunctions[router+".POST"] = worker
	} else {
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}

}
func (r *Runner) Get(router string, handelFunc func(ctx *Context), opts ...Option) {
	_, ok := r.handelFunctions[router]
	if !ok {
		worker := &Worker{
			Handel: []func(ctx *Context){handelFunc},
			Method: "GET",
			Path:   router,
			Config: &Config{},
		}
		if len(opts) > 0 {
			for _, opt := range opts {
				opt(worker.Config)
			}
		}
		r.handelFunctions[router+".GET"] = worker
	} else {
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}
}
