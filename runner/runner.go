package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/pkg/dto/runnerproject"
	"github.com/yunhanshu-net/sdk-go/env"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/request"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"time"
)

// New 创建一个新的Runner实例
func New() *Runner {
	runner, err := runnerproject.NewRunner(env.User, env.Name, env.Root, env.Version)
	if err != nil {
		panic(err)
	}
	return &Runner{
		idle:      5,
		detail:    runner,
		routerMap: make(map[string]*routerInfo),
		down:      make(chan struct{}, 1),
	}
}

// Runner 运行器结构体
type Runner struct {
	isDebug       bool
	detail        *runnerproject.Runner
	uuid          string
	idle          int64
	lastHandelTs  time.Time
	isClosed      bool
	natsConn      *nats.Conn
	natsSubscribe *nats.Subscription
	routerMap     map[string]*routerInfo
	down          chan struct{}
}

func (r *Runner) call(ctx context.Context, msg *nats.Msg) ([]byte, error) {

	data := msg.Data
	var req request.RunFunctionReq
	err1 := json.Unmarshal(data, &req)
	if err1 != nil {
		logger.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
	}

	runResponse, err1 := r.runFunction(ctx, &req)
	if err1 != nil {
		logger.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
	}
	marshal, err1 := json.Marshal(runResponse)
	if err1 != nil {
		logger.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
	}

	return marshal, nil
}

// connectNats 连接到NATS服务器并设置订阅
func (r *Runner) connectNats(ctx context.Context) error {
	now := time.Now()
	subject := r.detail.GetRequestSubject()

	// 增加NATS连接选项
	opts := []nats.Option{
		nats.Name(fmt.Sprintf("runner_%s_%s_%s", r.detail.User, r.detail.Name, r.uuid)),
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.WarnContextf(ctx, "NATS连接断开: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.InfoContextf(ctx, "NATS已重新连接")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			logger.ErrorContextf(ctx, "NATS错误: %v", err)
		}),
	}

	// 尝试连接NATS，带重试逻辑
	var connect *nats.Conn
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		logger.Infof("正在连接NATS服务器 (尝试: %d/%d)", i+1, maxRetries)
		connect, err = nats.Connect(nats.DefaultURL, opts...)
		if err == nil {
			break
		}
		logger.WarnContextf(ctx, "NATS连接失败，将在1秒后重试: %v", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf("无法连接NATS服务器: %w", err)
	}

	r.natsConn = connect

	// 设置消息处理
	subscribe, err := r.natsConn.QueueSubscribe(subject, subject, func(msg *nats.Msg) {
		r.lastHandelTs = time.Now()
		start := time.Now()

		// 创建响应消息
		respMsg := nats.NewMsg(msg.Reply)
		ctx1 := context.WithValue(context.Background(), constants.TraceID, msg.Header.Get(constants.TraceID))
		rspData, err := r.call(ctx1, msg)

		if err != nil {
			respMsg.Header.Set("code", "-1")
			respMsg.Header.Set("msg", err.Error())
			logger.Errorf("处理请求失败: %v", err)
		} else {
			respMsg.Data = rspData
			respMsg.Header.Set("code", "0")
		}

		// 响应请求
		if err := msg.RespondMsg(respMsg); err != nil {
			logger.ErrorContextf(ctx, "响应请求失败: %v", err)
			return
		}

		logger.DebugContextf(ctx, "请求处理完成，耗时: %v", time.Since(start))
	})
	logger.Infof("已连接到NATS服务器，监听主题: %s", subject)
	if err != nil {
		return fmt.Errorf("无法订阅主题 %s: %w", subject, err)
	}

	r.natsSubscribe = subscribe

	logger.InfoContextf(ctx, "uuid: %s", r.uuid)
	// 发送就绪消息
	msg := nats.NewMsg(r.uuid)
	msg.Header.Set("code", "0")

	//ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	respMsg, err := r.natsConn.RequestMsg(msg, time.Second*5)

	if err != nil {
		return fmt.Errorf("无法发送就绪消息: %w", err)
	}

	if respMsg.Header.Get("code") == "0" {
		logger.InfoContextf(ctx, "NATS连接成功，主题: %s，耗时: %v", subject, time.Since(now))
	} else {
		errMsg := respMsg.Header.Get("msg")
		return fmt.Errorf("NATS连接处理错误: %s", errMsg)
	}

	return nil
}

// close 安全关闭连接和订阅
func (r *Runner) close(ctx context.Context) error {
	// 防止重复关闭
	if r.isClosed {
		logger.DebugContextf(ctx, "Runner已经关闭，跳过")
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
			logger.DebugContextf(ctx, "清理订阅时出错: %v", err)
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
			ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 缩短超时时间
			defer cancel()

			// 尝试发送关闭通知，但不强制要求成功
			if msg, err := connToClose.RequestMsgWithContext(ctx1, newMsg); err != nil {
				logger.DebugContextf(ctx, "发送关闭通知失败: %v", err)
			} else if msg.Header.Get("code") != "0" {
				logger.DebugContextf(ctx, "关闭Runner返回错误: %s", msg.Header.Get("msg"))
			}
		}

		// 2.2 关闭连接
		connToClose.Close()
		logger.InfoContextf(ctx, "NATS连接已关闭")
	}

	logger.InfoContextf(ctx, "Runner资源清理完成，耗时: %v", time.Since(now))
	return closeErr
}
