package runner

type callback struct {
	Method string `json:"method"`
	Router string `json:"router"`
	Type   string `json:"type"`
}

const (
	callbackTypeOnCreated       = "onCreated"
	callbackTypeAfterDelete     = "afterDelete"
	callbackTypeOnVersionChange = "onVersionChange"
	callbackTypeBeforeClose     = "beforeClose"
	callbackTypeAfterClose      = "afterClose"
)

func (r *Runner) callback(ctx *HttpContext) error {
	var call callback
	err := ctx.Request.ShouldBindJSON(&call)
	if err != nil {
		return err
	}
	worker, exist := r.getRouterWorker(call.Router, call.Method)
	if exist {
		return nil
	}
	if worker.Config == nil {
		return nil
	}

	var callbackFunc func(ctx *HttpContext) error
	switch call.Type {
	case callbackTypeOnCreated:
		if worker.Config.OnCreated != nil {
			callbackFunc = worker.Config.OnCreated
		}
	case callbackTypeOnVersionChange:
		//遍历所有路由，只要有这个回调的，就执行
		if worker.Config.OnVersionChange != nil {
			callbackFunc = worker.Config.OnVersionChange
		}
	case callbackTypeAfterDelete:
		if worker.Config.AfterDelete != nil {
			callbackFunc = worker.Config.AfterDelete
		}
	case callbackTypeAfterClose:
		//遍历所有路由，只要有这个回调的，就执行
		if worker.Config.AfterClose != nil {
			callbackFunc = worker.Config.AfterClose
		}
	case callbackTypeBeforeClose:
		//遍历所有路由，只要有这个回调的，就执行
		if worker.Config.BeforeClose != nil {
			callbackFunc = worker.Config.BeforeClose
		}
	}
	if callbackFunc == nil {
		return nil
	}

	err = callbackFunc(ctx)
	if err != nil {
		return err
	}
	return nil

}
