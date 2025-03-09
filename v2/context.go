package v2

import (
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type Context struct {
	runner *model.Runner
	req    *request.RunnerRequest
	//transportConfig *TransportConfig
	about    int
	Request  *request.Request
	Response *response.Response
}

func (c *Context) About() {
	c.about = 1
}
