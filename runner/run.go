package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/request"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"runtime"
	"time"
)

func (r *Runner) connectCmd(cmd *cobra.Command, args []string) {
	//func (r *Runner) listen() {

	runnerId, err := cmd.Flags().GetString("runner_id")
	if err != nil {
		writeString(err.Error())
		return
	}
	ctx := context.Background()
	r.uuid = runnerId
	err = r.connectNats(ctx)
	if err != nil {
		writeString(err.Error())
		return
	} else {
		writeString("ok")
	}

	ticker := time.NewTicker(time.Second * 1)
	logger.Infof("listen uuid:%s\n", r.uuid)
	defer func() {
		ticker.Stop()
		// 使用统一的Shutdown函数而不是单独关闭资源
		Shutdown()
	}()

	for {
		select {
		case <-r.down:
			logger.Infof("%s runcher发起关闭请求，关闭连接", r.uuid)
			return
		case <-ticker.C:
			if r.idle > 0 {
				ts := time.Now().Unix()
				d := ts - r.lastHandelTs.Unix()
				if (ts - r.lastHandelTs.Unix()) > r.idle { //超过指定空闲时间的话需要释放进程
					logger.Infof(" %v没有处理消息，runner 自动关闭连接 idle config：%v", d, r.idle)
					return
				}
			}
		}
	}
	//}
}

func runCmd(cmd *cobra.Command, args []string) {
	// 从命令行获取参数值
	//method, err := cmd.Flags().GetString("method")
	//if err != nil {
	//	panic(err)
	//}
	//router, err := cmd.Flags().GetString("router")
	//if err != nil {
	//	panic(err)
	//}
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		panic(err)
	}

	traceID, err := cmd.Flags().GetString("trace_id")
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), constants.TraceID, traceID)
	var req request.RunFunctionReq
	err = jsonx.UnmarshalFromFile(file, &req)
	if err != nil {
		panic(err)
	}

	r.runCmd(ctx, &req)
}

// getRequest 从文件获取请求
func (r *Runner) getRequest(filePath string) (*request.RunFunctionReq, error) {
	// 检查文件是否存在
	//if _, err := os.Stat(filePath); os.IsNotExist(err) {
	//	return nil, fmt.Errorf("请求文件不存在: %s", filePath)
	//}

	var req request.RunFunctionReq
	//req.Request = new(request.Request)
	err := jsonx.UnmarshalFromFile(filePath, &req)
	if err != nil {
		return nil, fmt.Errorf("请求解析失败: %w", err)
	}
	return &req, nil
}

// getMemoryUsage 获取内存使用情况
func getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024)
}

// run 运行单次请求
func (r *Runner) runCmd(ctx context.Context, req *request.RunFunctionReq) {
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer func() {
	//	// 单次执行模式下，执行完请求后自动关闭资源
	//	if !r.isDebug {
	//		Shutdown()
	//	}
	//}()

	resp, err := r.runFunction(ctx, req)
	if err != nil {
		errorResp := &response.RunFunctionResp{
			Msg:  err.Error(),
			Code: -1,
			MetaData: map[string]interface{}{
				"error": err.Error(),
			},
		}
		marshal, _ := json.Marshal(errorResp)
		fmt.Println("<Response>" + string(marshal) + "</Response>")
		return
	}

	marshal, err := json.Marshal(resp)
	if err != nil {
		logger.Errorf("响应序列化失败: %s", err.Error())
		fmt.Println("<Response>{\"code\":500,\"msg\":\"响应序列化失败\"}</Response>")
		return
	}

	fmt.Println("<Response>" + string(marshal) + "</Response>")
}
