package runner

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	v2 "github.com/yunhanshu-net/sdk-go/model/response/v2"
	"time"
)

func (r *Runner) connect() error {
	msg := nats.NewMsg(fmt.Sprintf("runcher.%s.%s.%s.connect",
		r.detail.User, r.detail.Name, r.detail.Version))
	msg.Header.Set("connect", "req")
	msg.Header.Set("uuid", r.uuid)
	msg.Header.Set("subject", r.detail.GetSubject())

	connect, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logrus.Infof("connect:uuid failed: %s", r.uuid)
		return err
	}
	logrus.Info("connect:uuid success: " + r.uuid)
	r.conn = connect

	group := fmt.Sprintf("%s.%s.%s", r.detail.User, r.detail.Name, r.detail.Version)

	sub, err := connect.QueueSubscribe("runner.>", group, func(msg *nats.Msg) {
		r.lastHandelTs = time.Now()
		var reqMsg request.RunnerRequest
		err1 := sonic.Unmarshal(msg.Data, &reqMsg)
		if err1 != nil {
			panic(err1)
		}
		//ctx := &Context{req: &reqMsg, Request: reqMsg.Request, ResponseData: &response.ResponseData{}}
		httpContext := &HttpContext{
			Request:  reqMsg.Request,
			runner:   reqMsg.Runner,
			Response: &v2.ResponseData{},
		}
		err = r.runRequest(httpContext)
		marshal, err1 := sonic.Marshal(httpContext.Response)
		if err1 != nil {
			panic(err)
		}
		newMsg := nats.NewMsg(msg.Subject)
		newMsg.Data = marshal
		err1 = msg.RespondMsg(newMsg)
		if err1 != nil {
			panic(err1)
		}
	})
	if err != nil {
		panic(err)
	}
	r.sub = sub

	_, err = r.conn.RequestMsg(msg, time.Second*2)
	if err != nil {
		logrus.Infof("connect RequestMsg Ping uuid:%s err:%s", r.uuid, err)
		return err
	}

	logrus.Infof("connect done")
	return nil
}

func (r *Runner) close() error {
	msg := nats.NewMsg(fmt.Sprintf("runcher.%s.%s.%s.close",
		r.detail.User, r.detail.Name, r.detail.Version))
	msg.Header.Set("close", "req")
	msg.Header.Set("subject", r.detail.GetSubject())
	msg.Header.Set("uuid", r.uuid)
	_, err := r.conn.RequestMsg(msg, time.Second*2)
	if err != nil {
		return err
	}
	r.sub.Unsubscribe()
	r.conn.Close()
	return nil
}
