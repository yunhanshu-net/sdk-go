package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"time"
)

type natsConn struct {
	nats *nats.Conn
	sub  *nats.Subscription
}

func (n *natsConn) Connect(r *Runner) error {
	//todo 连接nats

	//调度引擎发起第一次握手
	connectKey := fmt.Sprintf("runner.connect.%s.%s", r.info.User, r.info.Soft)
	connectMsg := nats.NewMsg(connectKey)
	handelKey := fmt.Sprintf("runner.req.%s.%s", r.info.User, r.info.Soft)
	sub, err := n.nats.Subscribe(handelKey, r.onMsg)
	if err != nil { //连接失败响应
		connectMsg.Header.Set("status", "-1")
		connectMsg.Header.Set("msg", "connect_subscribe_failed")
		connectMsg.Header.Set("error", err.Error())
		_, err := n.nats.RequestMsg(connectMsg, time.Second*2) //连接失败
		if err != nil {
			return err
		}
		return err
	}
	n.sub = sub

	//连接成功响应
	connectMsg.Header.Set("status", "0")
	connectMsg.Header.Set("msg", "success")
	//发送连接报文响应给调度引擎，第二次握手
	rspConnectMsg, err := n.nats.RequestMsg(connectMsg, time.Second*3)
	if err != nil {
		return err
	}

	//确认调度引擎的回复信息，第三次握手
	if rspConnectMsg.Header.Get("status") != "0" { //说明失败
		return fmt.Errorf(rspConnectMsg.Header.Get("msg"))
	}
	r.conn = n
	//建立三次握手，连接成功，此时已经可以正常消费到消息了
	return nil
}

func (n *natsConn) Close() error {
	n.nats.Close()
	n.sub.Unsubscribe()
	return nil
}

func (n *natsConn) OnRequest(request *request.Request) error {

	return nil
}
func (n *natsConn) RequestChan(requestCh chan<- *request.Request) (err error) {
	return err
}
