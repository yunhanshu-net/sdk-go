package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/request"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"runtime"
)

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

	traceID, err := cmd.Flags().GetString("traceID")
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), constants.TraceID, traceID)
	var req requestv2.RunFunctionReq
	err = jsonx.UnmarshalFromFile(file, &req)
	if err != nil {
		panic(err)
	}

	r.runCmd(ctx, &req)
}

// getRequest 从文件获取请求
func (r *Runner) getRequest(filePath string) (*requestv2.RunFunctionReq, error) {
	// 检查文件是否存在
	//if _, err := os.Stat(filePath); os.IsNotExist(err) {
	//	return nil, fmt.Errorf("请求文件不存在: %s", filePath)
	//}

	var req requestv2.RunFunctionReq
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
func (r *Runner) runCmd(ctx context.Context, req *requestv2.RunFunctionReq) {
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
		logrus.Errorf("响应序列化失败: %s", err.Error())
		fmt.Println("<Response>{\"code\":500,\"msg\":\"响应序列化失败\"}</Response>")
		return
	}

	fmt.Println("<Response>" + string(marshal) + "</Response>")
}
