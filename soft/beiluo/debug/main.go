package main

import "github.com/yunhanshu-net/sdk-go/runner"

func main() {
	runner.Get("/hello", func(ctx *runner.HttpContext) {
		ctx.Response.JSON(map[string]interface{}{
			"hello": "world",
		}).Build()
	})
	//runner.Debug("beiluo", "debug", "v1", 30, "1211")
	runner.Run()
}
