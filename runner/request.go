package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/request"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"runtime"
	"runtime/debug"
	"time"
)

// runRequest 执行请求
func (r *Runner) runFunction(ctx context.Context, req *request.RunFunctionReq) (*response.RunFunctionResp, error) {
	router, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		routersJSON, _ := json.Marshal(r.routerMap)
		logger.ErrorContextf(ctx, "可用路由: %s", string(routersJSON))
		return nil, fmt.Errorf("路由未找到: [%s] %s", req.Method, req.Router)
	}

	// 使用defer-recover处理panic
	var result *response.RunFunctionResp
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				errMsg := fmt.Sprintf("请求处理panic: %v", r)
				logger.ErrorContextf(ctx, "%s\n调用栈: %s", errMsg, stack)
				err = fmt.Errorf(errMsg)
				// 这里打印是方便我出现错误时候可以直接在控制台看到日志
				fmt.Printf("err: %s\n调用栈: %s\n", errMsg, stack)
			}
		}()

		start := time.Now()
		var mStart runtime.MemStats
		var mEnd runtime.MemStats
		runtime.ReadMemStats(&mStart)
		_, rsp, callErr := router.call(ctx, req.Body)
		runtime.ReadMemStats(&mEnd)
		if callErr != nil {
			err = fmt.Errorf("路由调用失败: %w", callErr)
			return
		}

		// 记录执行时间
		elapsed := time.Since(start)
		if rsp.MetaData == nil {
			rsp.MetaData = make(map[string]interface{})
		}
		rsp.MetaData["cost"] = elapsed.String()
		rsp.MetaData["memory"] = getMemoryUsage()
		rsp.MetaData["cost_memory"] = fmt.Sprintf("%v", mEnd.Alloc-mStart.Alloc)

		result = rsp
	}()

	return result, err
}
