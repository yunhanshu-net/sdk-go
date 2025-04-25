package runner

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

// connectNats 连接到NATS服务器并设置订阅
func (r *Runner) connectNats() error {
	now := time.Now()
	subject := r.detail.GetRequestSubject()

	// 增加NATS连接选项
	opts := []nats.Option{
		nats.Name(fmt.Sprintf("runner_%s_%s_%s", r.detail.User, r.detail.Name, r.uuid)),
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logrus.Warnf("NATS连接断开: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logrus.Info("NATS已重新连接")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			logrus.Errorf("NATS错误: %v", err)
		}),
	}

	// 尝试连接NATS，带重试逻辑
	var connect *nats.Conn
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		logrus.Infof("正在连接NATS服务器 (尝试: %d/%d)", i+1, maxRetries)
		connect, err = nats.Connect(nats.DefaultURL, opts...)
		if err == nil {
			break
		}
		logrus.Warnf("NATS连接失败，将在1秒后重试: %v", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf("无法连接NATS服务器: %w", err)
	}

	r.natsConn = connect
	logrus.Infof("已连接到NATS服务器，监听主题: %s", subject)

	// 设置消息处理
	subscribe, err := r.natsConn.QueueSubscribe(subject, subject, func(msg *nats.Msg) {
		r.lastHandelTs = time.Now()
		start := time.Now()

		// 创建响应消息
		respMsg := nats.NewMsg(msg.Reply)
		rspData, err := r.call(msg)

		if err != nil {
			respMsg.Header.Set("code", "-1")
			respMsg.Header.Set("msg", err.Error())
			logrus.Errorf("处理请求失败: %v", err)
		} else {
			respMsg.Data = rspData
			respMsg.Header.Set("code", "0")
		}

		// 响应请求
		if err := msg.RespondMsg(respMsg); err != nil {
			logrus.Errorf("响应请求失败: %v", err)
		}

		logrus.Debugf("请求处理完成，耗时: %v", time.Since(start))
	})

	if err != nil {
		return fmt.Errorf("无法订阅主题 %s: %w", subject, err)
	}

	r.natsSubscribe = subscribe

	// 发送就绪消息
	msg := nats.NewMsg(r.uuid)
	msg.Header.Set("code", "0")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	respMsg, err := r.natsConn.RequestMsgWithContext(ctx, msg)

	if err != nil {
		return fmt.Errorf("无法发送就绪消息: %w", err)
	}

	if respMsg.Header.Get("code") == "0" {
		logrus.Infof("NATS连接成功，主题: %s，耗时: %v", subject, time.Since(now))
	} else {
		errMsg := respMsg.Header.Get("msg")
		return fmt.Errorf("NATS连接处理错误: %s", errMsg)
	}

	return nil
}

// close 安全关闭连接和订阅
func (r *Runner) close() error {
	// 防止重复关闭
	if r.isClosed {
		logrus.Debug("Runner已经关闭，跳过")
		return nil
	}

	// 标记为已关闭
	r.isClosed = true
	now := time.Now()

	var closeErr error

	// 1. 先尝试清理订阅
	if r.natsSubscribe != nil {
		subToClose := r.natsSubscribe
		r.natsSubscribe = nil // 立即置空，防止重复关闭

		if err := subToClose.Drain(); err != nil {
			logrus.Warnf("清理订阅时出错: %v", err)
			closeErr = fmt.Errorf("清理订阅错误: %w", err)
			// 继续执行，不中断关闭流程
		}
	}

	// 2. 处理NATS连接
	if r.natsConn != nil {
		connToClose := r.natsConn
		r.natsConn = nil // 立即置空，防止重复关闭

		// 发送关闭通知（尽最大努力）
		if r.detail != nil { // 检查detail是否为nil
			// 2.1 尝试发送关闭通知
			subject := "close.runner"
			newMsg := nats.NewMsg(subject)
			newMsg.Header.Set("version", r.detail.Version)
			newMsg.Header.Set("user", r.detail.User)
			newMsg.Header.Set("name", r.detail.Name)
			newMsg.Data = []byte(r.uuid)

			// 设置请求超时
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 缩短超时时间
			defer cancel()

			// 尝试发送关闭通知，但不强制要求成功
			if msg, err := connToClose.RequestMsgWithContext(ctx, newMsg); err != nil {
				logrus.Warnf("发送关闭通知失败: %v", err)
			} else if msg.Header.Get("code") != "0" {
				logrus.Warnf("关闭Runner返回错误: %s", msg.Header.Get("msg"))
			}
		}

		// 2.2 关闭连接
		connToClose.Close()
		logrus.Info("NATS连接已关闭")
	}

	logrus.Infof("Runner资源清理完成，耗时: %v", time.Since(now))
	return closeErr
}
