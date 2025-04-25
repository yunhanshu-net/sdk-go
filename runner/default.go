package runner

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var r = New()
var shutdownOnce sync.Once // 确保只关闭一次

// Shutdown 统一的资源关闭入口，处理所有资源的释放
func Shutdown() {
	shutdownOnce.Do(func() {
		logrus.Info("开始执行系统关闭...")

		// 1. 先关闭Runner连接，包括NATS连接等
		if err := r.close(); err != nil {
			logrus.Errorf("关闭Runner连接失败: %v", err)
		}

		// 2. 关闭所有数据库连接
		CloseAllDBs()

		// 3. 这里添加其他需要关闭的资源
		// ...

		logrus.Info("系统关闭完成")
	})
}

// SetupSignalHandler 设置信号处理，捕获SIGINT、SIGTERM信号
func SetupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		logrus.Infof("收到信号: %v, 开始优雅退出...", sig)
		Shutdown()
		os.Exit(0)
	}()
}

func Run() error {
	// 先设置信号处理
	SetupSignalHandler()

	// 执行原有逻辑，但不再需要在这里处理资源关闭
	return r.Run()
}

func Debug(user, runner, version string, idle int64, uuid string) error {
	// 先设置信号处理
	SetupSignalHandler()

	// 执行原有逻辑，但不再需要在这里处理资源关闭
	return r.Debug(user, runner, version, idle, uuid)
}

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*ApiConfig) {
	r.get(router, handler, config...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*ApiConfig) {
	r.post(router, handler, config...)
}

func runHandel(method, router string, body string) error {
	key := fmtKey(router, method)
	info := r.routerMap[key]
	if info == nil {
		return fmt.Errorf("router not found")
	}
	req, rsp, err := info.call(context.Background(), body)
	if err != nil {
		logrus.Errorf("路由调用失败: %v", err)
		return err
	}
	if req == nil || rsp == nil {
		err = fmt.Errorf("返回结果为空")
		logrus.Errorf("路由返回结果为空: %s", key)
		return err
	}
	return nil
}
