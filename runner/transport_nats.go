package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type transportNats struct {
	natsConn *nats.Conn
	natsSub  *nats.Subscription
}

func newTransportNats(ctxInfo *ContextInfo) (*transportNats, error) {
	url := ctxInfo.Metadata["nats-url"]
	if url == "" {
		url = nats.DefaultURL
	}
	group := ctxInfo.Metadata["nats-group"]
	if group == "" {
		group = fmt.Sprintf("%s.%s.%s", ctxInfo.User, ctxInfo.Soft, ctxInfo.Version)
	}

	transport := &transportNats{}
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	transport.natsConn = conn

	//函数请求
	//runner.user.soft.version.run
	//header 携带路由和
	//函数请求

	//关闭连接
	//runner.user.soft.version.close 关闭连接请求
	//关闭连接

	//心跳检测，探针，判断调度引擎是否还存活正常
	//runner.user.soft.version.heartbeat_check
	//心跳检测，探针

	subject := fmt.Sprintf("runner.>")
	conn.QueueSubscribe(subject, group, func(msg *nats.Msg) {

	})

}

func (t *transportNats) ReadMessage() (*TransportMsg, error) {
	//TODO implement me
	panic("implement me")
}

func (t *transportNats) WriteMessage(msg *TransportMsg) error {
	//TODO implement me
	panic("implement me")
}

func (t *transportNats) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (t *transportNats) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (t *transportNats) Close() error {
	//TODO implement me
	panic("implement me")
}
