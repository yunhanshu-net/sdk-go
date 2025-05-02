package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/constants"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
)

// New 创建一个新的Runner实例
func New() *Runner {
	return &Runner{
		idle:      5,
		routerMap: make(map[string]*routerInfo),
		down:      make(chan struct{}, 1),
	}
}

// Runner 运行器结构体
type Runner struct {
	isDebug       bool
	detail        *model.Runner
	uuid          string
	args          []string
	idle          int64
	lastHandelTs  time.Time
	isClosed      bool
	natsConn      *nats.Conn
	natsSubscribe *nats.Subscription
	routerMap     map[string]*routerInfo
	down          chan struct{}
}

// init 初始化Runner
func (r *Runner) init(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("参数不足，至少需要3个参数")
	}

	r.args = args
	r.detail = &model.Runner{}

	// 解析Runner名称
	split := strings.Split(r.args[0], "_")
	if len(split) > 1 {
		r.detail.User = strings.ReplaceAll(split[0], "./", "")
		r.detail.Name = split[1]
	} else {
		return fmt.Errorf("runner名称格式不正确: %s", r.args[0])
	}

	logrus.Infof("Runner详情: %+v", r.detail)

	// 设置最大处理器数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 注册内置路由
	r.registerBuiltInRouters()

	// 获取请求
	req, err := r.getRequest(r.args[2])
	if err != nil {
		return fmt.Errorf("获取请求失败: %w", err)
	}
	var ctx context.Context
	if req.Request != nil {
		ctx = context.WithValue(context.Background(), constants.TraceID, req.Request.TraceID)
	} else {
		ctx = context.Background()
	}

	if req != nil {
		r.detail = req.Runner
		if r.uuid == "" {
			r.uuid = req.UUID
		}
	}

	// 处理连接模式
	if r.args[1] == "_connect" {
		if req.TransportConfig != nil && req.TransportConfig.IdleTime != 0 {
			r.idle = int64(req.TransportConfig.IdleTime)
		}

		// 异步连接NATS
		errChan := make(chan error, 1)
		go func() {
			err := r.connectNats(ctx)
			if err != nil {
				errChan <- err
			}
		}()

		// 等待连接结果
		select {
		case err := <-errChan:
			return fmt.Errorf("NATS连接失败: %w", err)
		case <-time.After(5 * time.Second):
			// 连接成功或超时
		}

		r.listen()
		logrus.Infof("UUID: %s 监听已停止", r.uuid)
		return nil
	}

	// 单次执行模式
	r.run(ctx, req)
	return nil
}

// registerBuiltInRouters 注册内置路由
func (r *Runner) registerBuiltInRouters() {
	r.get("/_env", env)
	r.get("/_ping", ping)
	r.get("/_getApiInfos", r._getApiInfos)
	r.get("/_getApiInfo", r._getApiInfo)
	r.post("/_callback", r._callback)
	r.post("/_sysCallback/sysOnVersionChange", r._sysOnVersionChange)
}

// getRequest 从文件获取请求
func (r *Runner) getRequest(filePath string) (*request.RunnerRequest, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("请求文件不存在: %s", filePath)
	}

	var req request.RunnerRequest
	req.Request = new(request.Request)
	err := jsonx.UnmarshalFromFile(filePath, &req)
	if err != nil {
		return nil, fmt.Errorf("请求解析失败: %w", err)
	}
	return &req, nil
}

// getRouter 获取路由
func (r *Runner) getRouter(router string, method string) (worker *routerInfo, exist bool) {
	worker, ok := r.routerMap[fmtKey(router, method)]
	return worker, ok
}

// runRequest 执行请求
func (r *Runner) runRequest(ctx context.Context, req *request.Request) (*response.Data, error) {
	router, exist := r.getRouter(req.Route, req.Method)
	if !exist {
		routersJSON, _ := json.Marshal(r.routerMap)
		logger.ErrorContextf(ctx, "可用路由: %s", string(routersJSON))
		return nil, fmt.Errorf("路由未找到: [%s] %s", req.Method, req.Route)
	}

	// 使用defer-recover处理panic
	var result *response.Data
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
		_, rsp, callErr := router.call(ctx, req.Body)
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

		result = rsp
	}()

	return result, err
}

// getMemoryUsage 获取内存使用情况
func getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024)
}

// run 运行单次请求
func (r *Runner) run(ctx context.Context, req *request.RunnerRequest) {
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		// 单次执行模式下，执行完请求后自动关闭资源
		if !r.isDebug {
			Shutdown()
		}
	}()

	resp, err := r.runRequest(ctx, req.Request)
	if err != nil {
		errorResp := &response.Data{
			Msg: "请求处理失败",
			MetaData: map[string]interface{}{
				"error": err.Error(),
			},
		}
		marshal, _ := sonic.Marshal(errorResp)
		fmt.Println("<Response>" + string(marshal) + "</Response>")
		return
	}

	marshal, err := sonic.Marshal(resp)
	if err != nil {
		logrus.Errorf("响应序列化失败: %s", err.Error())
		fmt.Println("<Response>{\"code\":500,\"msg\":\"响应序列化失败\"}</Response>")
		return
	}

	fmt.Println("<Response>" + string(marshal) + "</Response>")
}

// Run 运行Runner
//func (r *Runner) Run() error {
//	// 先初始化
//	err := r.init(os.Args)
//	if err != nil {
//		return err
//	}
//
//	// 如果不是连接模式，而是单次执行模式，这里已经执行完毕
//	// 我们应该在这里释放资源，而不是依赖run方法
//	if r.args[1] != "_connect" {
//		// 在单次执行模式下，runner.go中的run方法已经处理了请求
//		// 不需要再次调用Shutdown，避免重复关闭
//	} else {
//		// 在连接模式下，资源会在listen中通过信号处理或超时机制关闭
//	}
//
//	return nil
//}
