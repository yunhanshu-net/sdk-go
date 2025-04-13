package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/response"
	"testing"
)

type HelloReq struct {
	Msg string `json:"msg"`
}
type HelloResp struct {
	Rsp string `json:"rsp"`
}

func Hello(ctx *Context, req *HelloReq, resp response.Response) error {
	return resp.JSON(&HelloResp{Rsp: req.Msg + ":rsp"}).Build()
}

func TestName(t *testing.T) {
	Get("/hello", Hello)
	s := `{"msg":"hello"}`

	err := runHandel("GET", "/hello", s)
	if err != nil {
		panic(err)
	}
}

func BenchmarkName(b *testing.B) {
	Get("/hello", Hello)
	s := `{"msg":"hello"}`
	for i := 0; i < b.N; i++ {
		err := runHandel("GET", "/hello", s)
		if err != nil {
			panic(err)
		}
	}
}
