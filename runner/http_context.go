package runner

import (
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	v2 "github.com/yunhanshu-net/sdk-go/model/response/v2"
)

type HttpContext struct {
	//runner *Context
	//req    *request.RunnerRequest
	Request  *request.Request
	Response v2.Response
	runner   *model.Runner
}
