package main

import (
	"github.com/yunhanshu-net/sdk-go/runner"
)

func main() {

	runner.Get("/hello", func(ctx *runner.Context) {
		mp := make(map[string]interface{})
		mp["hello"] = "world"
		ctx.Response.OKWithJSON(mp)
	})

	err := runner.Run()
	if err != nil {
		panic(err)
	}

	//r.Debug("beiluo", "apphub", "v1", 0, "uuid1111")
	//r.Debug("beiluo", "apphub", "v1", 0, "uuid1111")
}
