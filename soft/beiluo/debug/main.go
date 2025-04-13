package main

import (
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"github.com/yunhanshu-net/sdk-go/soft/beiluo/debug/version/v1/api/test"
)

type HelloResp struct {
	Hello string `json:"hello"`
	World string
}
type HelloReq struct {
}

func main() {
	test.Setup()
	runner.Get("/hello", func(ctx *runner.Context, req *HelloReq, resp response.Response) error {
		return resp.JSON(HelloResp{Hello: "hello 12", World: "World 121"}).Build()
	})
	//runner.Debug("beiluo", "debug", "v1", 30, "1211")
	runner.Run()
}
