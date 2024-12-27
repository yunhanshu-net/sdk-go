package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/logger"
)

func (c *Context) GetLogger() *logger.Logger {

	mp := make(map[string]interface{})
	mp["a_tenant"] = c.runner.info.User
	mp["a_soft"] = c.runner.info.Soft
	mp["a_command"] = c.runner.info.Command
	if c.runner != nil {
		if c.runner.info.Version != "" {
			mp["a_version"] = c.runner.info.Version
		}
	}

	mp["a_soft_classify"] = fmt.Sprintf("/%s/%s", c.runner.info.User, c.runner.info.Soft)
	return &logger.Logger{DataMap: mp}
}
