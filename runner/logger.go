package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/logger"
)

func (c *Context) GetLogger() *logger.Logger {

	mp := make(map[string]interface{})
	mp["a_tenant"] = c.transportConfig.User
	mp["a_trace_id"] = c.Request.TraceID
	mp["a_soft"] = c.transportConfig.Runner
	mp["a_command"] = c.transportConfig.Route
	if c.runner != nil {
		if c.transportConfig.Version != "" {
			mp["a_version"] = c.transportConfig.Version
		}
	}

	mp["a_soft_classify"] = fmt.Sprintf("/%s/%s", c.transportConfig.User, c.transportConfig.Runner)
	return &logger.Logger{DataMap: mp}
}
