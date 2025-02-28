package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"strings"
	"sync"
	"time"
)

type transportNats struct {
	wg               *sync.WaitGroup
	readMsgCount     int
	responseMsgCount int
	natsConn         *nats.Conn
	natsSub          *nats.Subscription
	msgList          chan *TransportMsg
	contextInfo      *TransportConfig
}

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

func newTransportNats(transportConfig *TransportConfig) (trs *transportNats, err error) {
	url := transportConfig.Metadata["nats-url"]
	if url == "" {
		url = nats.DefaultURL
	}
	group := transportConfig.Metadata["nats-group"]
	if group == "" {
		group = fmt.Sprintf("%s.%s.%s", transportConfig.User, transportConfig.Runner, transportConfig.Version)
	}

	transport := &transportNats{
		msgList:     make(chan *TransportMsg, 2000),
		contextInfo: transportConfig,
	}
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	err = transport.Ping() //建立连接
	if err != nil {
		return nil, err
	}

	transport.natsConn = conn
	subject := fmt.Sprintf("runner.%s.%s.%s.*", transportConfig.User, transportConfig.Runner, transportConfig.Version)
	sub, err := conn.QueueSubscribe(subject, group, func(msg *nats.Msg) {
		transport.wg.Add(1)
		transport.readMsgCount++
		subjects := strings.Split(msg.Subject, ".")
		cmd := subjects[len(subjects)-1]

		fmt.Println("receive:", string(msg.Data))
		headers := make(map[string][]string)
		for k, v := range headers {
			headers[k] = v
		}
		if cmd == MsgTypeRun {
			trMsg := &TransportMsg{
				msg:       msg,
				Data:      msg.Data,
				Headers:   headers,
				Subject:   msg.Subject,
				transport: transport,
			}
			transport.msgList <- trMsg
		}
	})
	if err != nil {
		return nil, err
	}
	trs.natsSub = sub
	return transport, nil

}

func (t *transportNats) ReadMessage() <-chan *TransportMsg {
	return t.msgList
}

func (t *transportNats) Ping() error {
	msg := nats.NewMsg(fmt.Sprintf("runner.%s.%s.%s.connect",
		t.contextInfo.User, t.contextInfo.Runner, t.contextInfo.Version))
	msg.Header.Set("connect", "req")
	requestMsg, err := t.natsConn.RequestMsg(msg, time.Second*2)
	if err != nil {
		return err
	}
	res := requestMsg.Header.Get("connect")
	if res == "resp" { //说明关闭成功
		return nil
	}
	return fmt.Errorf(requestMsg.Header.Get("msg"))
}

func (t *transportNats) Close() error {
	//先发送关闭请求
	msg := nats.NewMsg(fmt.Sprintf("runner.%s.%s.%s.close",
		t.contextInfo.User, t.contextInfo.Runner, t.contextInfo.Version))
	msg.Header.Set("close", "req")
	requestMsg, err := t.natsConn.RequestMsg(msg, time.Second*2)
	if err != nil {
		return err
	}
	res := requestMsg.Header.Get("close")

	if res == "resp" { //说明关闭成功
		t.natsSub.Unsubscribe()
		t.natsConn.Close()

		return nil
	}
	return fmt.Errorf(requestMsg.Header.Get("msg"))
}

func (t *transportNats) Wait() <-chan struct{} {
back:
	t.wg.Wait()
	if t.readMsgCount != t.responseMsgCount {
		goto back
	}

	return make(<-chan struct{})
}
