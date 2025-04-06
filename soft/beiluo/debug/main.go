package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/runner"
)

type Hello struct {
	Hello string `json:"hello"`
	World string
}

func main() {
	defer func() {
		logrus.Infof("done")
	}()
	runner.Get("/hello", func(ctx *runner.HttpContext) {
		ctx.Response.JSON(Hello{
			Hello: "hello",
			World: "World",
		}).Build()
	})
	//runner.Debug("beiluo", "debug", "v1", 30, "1211")
	runner.Run()
}
