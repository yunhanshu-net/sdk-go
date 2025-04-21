package runner

import (
	"context"
)

type Context struct {
	context.Context
}

func (c *Context) GetUsername() string {
	return ""
}
