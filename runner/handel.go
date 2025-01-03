package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type Worker struct {
	Handel []func(ctx *Context)
	Path   string
	Method string
	Config *Config
}

func (r *Runner) handel(io IO) error {
	req, err := io.Input(r)
	if err != nil {
		return err
	}
	work := r.handelFunctions[req.Url+"."+req.Method]
	ctx := &Context{
		runner:   r,
		Request:  req,
		Response: &response.Response{},
	}
	if work == nil {
		if r.notFound != nil {
			r.notFound(ctx)
			io.Output(r, ctx.Response)
			return nil
		} else {
			return fmt.Errorf("no work found")
		}
	}
	for _, handel := range work.Handel {
		if ctx.about == 1 {
			break
		}
		//todo 这里执行用户注册的路由函数
		handel(ctx)
	}

	io.Output(r, ctx.Response)
	return nil
}

func (r *Runner) handelRequest(req *request.Request) (*Context, error) {
	work := r.handelFunctions[req.Url+"."+req.Method]
	ctx := &Context{
		runner:   r,
		Request:  req,
		Response: &response.Response{},
	}
	if work == nil {
		if r.notFound != nil {
			r.notFound(ctx)
			//io.Output(r, ctx.Response)
			r.conn.Response(ctx.Response)
			return ctx, nil
		} else {
			return nil, fmt.Errorf("no work found")
		}
	}
	for _, handel := range work.Handel {
		if ctx.about == 1 {
			break
		}
		//todo 这里执行用户注册的路由函数
		handel(ctx)
	}
	r.conn.Response(ctx.Response)
	//io.Output(r, ctx.Response)
	return ctx, nil
}
