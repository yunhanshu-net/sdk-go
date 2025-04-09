package runner

import (
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type HttpContext struct {
	Request  *request.Request
	Response response.Response
	runner   *model.Runner
}

func (c *HttpContext) GetUser() string {
	return ""
}
