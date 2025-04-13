package runner

import (
	"context"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

var r = New()

func Run() error {
	return r.Run()
}
func Debug(user, runner, version string, idle int64, uuid string) error {
	return r.Debug(user, runner, version, idle, uuid)
}

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*ApiConfig) {
	r.get(router, handler, config...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*ApiConfig) {
	r.post(router, handler, config...)
}

func runHandel(method, router string, body string) error {
	key := fmtKey(router, method)
	info := r.routerMap[key]
	if info == nil {
		return fmt.Errorf("router not found")
	}
	req, rsp, err := info.call(context.Background(), body)
	if err != nil {
		panic(err)
	}
	if req == nil || rsp == nil {
		panic(err)
	}
	return nil
}
