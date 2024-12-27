package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

type Context struct {
	runner   *Runner
	about    int
	Request  *request.Request
	Response *response.Response
}

func (c *Context) About() {
	c.about = 1
}
