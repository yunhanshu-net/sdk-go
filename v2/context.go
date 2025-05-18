package v2

import (
	"context"
	"github.com/yunhanshu-net/pkg/constants"
)

type Context struct {
	context.Context
}

func (c *Context) getTraceId() string {
	value := c.Context.Value(constants.TraceID)
	if value == nil {
		return ""
	}
	v, ok := value.(string)
	if ok {
		return v
	}
	return ""
}
