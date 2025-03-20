package runner

import "fmt"

func (r *Runner) Post(router string, handelFunc func(ctx *Context), config ...*Config) {
	_, ok := r.handelFunctions[router]
	if !ok {
		worker := &Worker{
			Handel: []func(*Context){handelFunc},
			Method: "POST",
			Path:   router,
			Config: &Config{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}
		r.handelFunctions[router+".POST"] = worker
	} else {
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}

}
func (r *Runner) Get(router string, handelFunc func(ctx *Context), config ...*Config) {
	_, ok := r.handelFunctions[router]
	if !ok {
		fmt.Println("get---------- !ok")
		worker := &Worker{
			Handel: []func(ctx *Context){handelFunc},
			Method: "GET",
			Path:   router,
			Config: &Config{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.Config = config[0]
		}
		r.handelFunctions[router+".GET"] = worker
	} else {
		fmt.Println("get---------- ok")
		r.handelFunctions[router].Handel = append(r.handelFunctions[router].Handel, handelFunc)
	}
}
