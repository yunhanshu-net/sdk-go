package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/logger"
)

func (c *Context) GetLogger() *logger.Logger {

	mp := make(map[string]interface{})
	mp["a_tenant"] = c.runner.User
	mp["a_soft"] = c.runner.Soft
	mp["a_command"] = c.runner.Command
	if c.runner != nil {
		if c.runner.Version != "" {
			mp["a_version"] = c.runner.Version
		}
	}

	mp["a_soft_classify"] = fmt.Sprintf("/%s/%s", c.runner.User, c.runner.Soft)
	return &logger.Logger{
		DataMap: mp,
	}
}
