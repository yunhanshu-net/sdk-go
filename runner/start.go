package runner

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yunhanshu-net/pkg/x/jsonx"
	"github.com/yunhanshu-net/sdk-go/env"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func writeJSON(el interface{}) {
	fmt.Println("<Response>" + jsonx.String(el) + "</Response>")
}
func writeString(msg string) {
	fmt.Println("<Response>" + msg + "</Response>")
}

var cmd *cobra.Command
var r *Runner

func initRunner() {
	if r == nil {
		r = New()
	}
	if cmd == nil {
		cmd = newCmd()
	}
}

func newCmd() *cobra.Command {
	app := &cobra.Command{Use: fmt.Sprintf("%s_%s_%s", env.User, env.Name, env.Version), Short: "根命令"}
	run := &cobra.Command{Use: "run", Short: "运行函数", Run: runCmd}
	// 添加标志定义到 run 子命令
	run.Flags().String("file", "", "JSON 请求文件路径")   // 长格式 --file
	run.Flags().String("method", "POST", "HTTP 方法") // 长格式 --method
	run.Flags().String("router", "", "请求路由路径")      // 长格式 --router
	run.Flags().String("trace_id", "", "请求跟踪ID")    // 长格式 --trace_id
	connect := &cobra.Command{Use: "connect", Short: "建立连接", Run: r.connectCmd}
	connect.Flags().String("runner_id", "", "runnerId") // 长格式 --connect_id
	//syscall := &cobra.Command{Use: "syscall", Short: "系统回调", Run: r.syscallCmd}
	userCall := &cobra.Command{Use: "usercall", Short: "用户回调", Run: r.userCallCmd}
	userCall.Flags().String("type", "", "回调类型")
	userCall.Flags().String("method", "", "回调方法")
	userCall.Flags().String("router", "", "回调路由")
	userCall.Flags().String("file", "", "回调body文件路径")
	userCall.Flags().String("trace_id", "", "请求跟踪ID") // 长格式 --trace_id

	apis := &cobra.Command{Use: "apis", Short: "apis", Run: r.apisCmd}

	app.AddCommand(run)
	app.AddCommand(apis)
	app.AddCommand(connect)
	app.AddCommand(userCall)
	return app
}

var shutdownOnce sync.Once // 确保只关闭一次

func Run() error {
	//SetupSignalHandler()
	initRunner()
	r.registerBuiltInRouters()
	if err := cmd.Execute(); err != nil {
		logger.Fatal(err.Error())

	}
	Shutdown()
	return nil
}

// Shutdown 统一的资源关闭入口，处理所有资源的释放
func Shutdown() {
	shutdownOnce.Do(func() {
		logger.Info("开始执行系统关闭...")

		// 1. 先关闭Runner连接，包括NATS连接等
		if err := r.close(context.Background()); err != nil {
			logger.Errorf("关闭Runner连接失败: %v", err)
		}

		// 2. 关闭所有数据库连接
		CloseAllDBs()

		// 3. 这里添加其他需要关闭的资源
		// ...

		logger.Info("系统关闭完成")
	})
}

// SetupSignalHandler 设置信号处理，捕获SIGINT、SIGTERM信号
func SetupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		logger.Infof("收到信号: %v, 开始优雅退出...", sig)
		Shutdown()
		os.Exit(0)
	}()
}
