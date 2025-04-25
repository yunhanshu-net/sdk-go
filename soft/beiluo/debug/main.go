package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"github.com/yunhanshu-net/sdk-go/soft/beiluo/debug/version/v1/api/apiinfo"
	"github.com/yunhanshu-net/sdk-go/soft/beiluo/debug/version/v1/api/calc"
)

type HelloResp struct {
	Hello string `json:"hello"`
	World string
}
type HelloReq struct {
}

func main() {
	// 初始化API模块
	calc.Setup()
	apiinfo.Setup()

	// 注册路由
	runner.Get("/hello", func(ctx *runner.Context, req *HelloReq, resp response.Response) error {
		return resp.JSON(HelloResp{Hello: "hello", World: "World"}).Build()
	})

	// 运行应用 - Shutdown处理已经在runner库中实现
	if err := runner.Run(); err != nil {
		logrus.Fatalf("应用启动失败: %v", err)
	}
}
