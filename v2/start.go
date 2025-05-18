package response

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var cmd = newCmd()
var r = New()

func newCmd() *cobra.Command {
	app := &cobra.Command{
		Use:   "app",
		Short: "根命令",
	}
	run := &cobra.Command{
		Use:   "run",
		Short: "运行函数",
		Run:   runCmd,
	}

	ping := &cobra.Command{
		Use:   "ping",
		Short: "ping",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("pong")
		},
	}
	app.AddCommand(run)
	app.AddCommand(ping)
	return app
}

var shutdownOnce sync.Once // 确保只关闭一次

func Run() error {
	//SetupSignalHandler()
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
