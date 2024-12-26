package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type Context struct {
	runner   *Runner
	Request  *request.Request
	Response *response.Response
}
