package array

import (
	"github.com/yunhanshu-net/sdk-go/runner"
)

func init() {
	runner.Get("/array/diff", Diff)
}

type DiffResp struct {
	Arr string `json:"arr"`
}

func Diff(ctx *runner.HttpContext) {

	err := ctx.Response.JSON(&DiffResp{Arr: "a,b,c,d"}).Build()
	if err != nil {
		panic(err)
	}
}
