package runner

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/model/status"
	"sync"
	"time"
)

type natsConn struct {
	info      *Info
	reqWg     *sync.WaitGroup
	nats      *nats.Conn
	sub       *nats.Subscription
	requestCh chan *request.Request
}

func (n *natsConn) Connect() error {
	//todo 连接nats

	//调度引擎发起第一次握手
	connectKey := fmt.Sprintf("runner.connect.%s.%s", n.info.User, n.info.Soft)
	connectMsg := nats.NewMsg(connectKey)
	handelKey := fmt.Sprintf("runner.req.%s.%s", n.info.User, n.info.Soft)
	sub, err := n.nats.Subscribe(handelKey, func(msg *nats.Msg) {
		//r.wg.Add(1)
		//defer r.wg.Done()

		n.reqWg.Add(1)
		requestMsg := &request.Request{}
		err := json.Unmarshal(msg.Data, requestMsg)
		if err != nil {
			fmt.Println(err)
			defer n.reqWg.Done()
		}
		n.requestCh <- requestMsg
		//err := r.conn.OnRequest(req)
		//if err != nil {
		//	panic(err)
		//}
	})
	if err != nil { //连接失败响应
		connectMsg.Header.Set("status", fmt.Sprintf("%d", status.ConnectError.Status))
		connectMsg.Header.Set("msg", status.ConnectError.Message)
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
	//r.conn = n
	//建立三次握手，连接成功，此时已经可以正常消费到消息了
	return nil
}

func (n *natsConn) Close() error {
	defer n.sub.Unsubscribe()
	defer n.nats.Close()
	closeKey := fmt.Sprintf("runner.close.%s.%s", n.info.User, n.info.Soft)
	msg := nats.NewMsg(closeKey)
	rsp, err := n.nats.RequestMsg(msg, time.Second*3)
	if err != nil {
		return err
	}
	if rsp.Header.Get("status") != "0" {
		return fmt.Errorf(rsp.Header.Get("msg"))
	}

	return nil
}

func (n *natsConn) Response(resp *response.Response) error {
	respKey := fmt.Sprintf("runner.resp.%s.%s", n.info.User, n.info.Soft)
	msg := nats.NewMsg(respKey)
	err := n.nats.PublishMsg(msg)
	return err
}
