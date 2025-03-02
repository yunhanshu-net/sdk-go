package main

import "github.com/yunhanshu-net/sdk-go/runner"

func main() {

	r := runner.New()
	r.Get("/hello", func(ctx *runner.Context) {
		mp := make(map[string]interface{})
		mp["hello"] = "world"
		ctx.Response.OKWithJSON(mp)
	})

	//r.Run()
	r.Debug("beiluo", "apphub", "v1", 0)
}
