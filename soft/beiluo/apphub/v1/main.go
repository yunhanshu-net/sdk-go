package main

import (
	"github.com/yunhanshu-net/sdk-go/runner"
	_ "github.com/yunhanshu-net/sdk-go/soft/beiluo/apphub/v1/biz/user"
)

func main() {

	runner.Get("/hello", func(ctx *runner.HttpContext) {
		mp := make(map[string]interface{})
		mp["hello"] = "world"
		ctx.Response.JSON(mp).Build()
	})
	//runner.Get("user/list", user.List)

	err := runner.Run()
	if err != nil {
		panic(err)
	}

	//r.Debug("beiluo", "apphub", "v1", 0, "uuid1111")
	//r.Debug("beiluo", "apphub", "v1", 0, "uuid1111")
}
