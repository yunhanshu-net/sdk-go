package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

func (r *Runner) connectNats() error {
	now := time.Now()
	subject := r.detail.GetRequestSubject()
	connect, err := nats.Connect(nats.DefaultURL)
	logrus.Infof("subject:%s", subject)
	if err != nil {
		logrus.Errorf("connectNats nats.Connect err:%s", err.Error())
		return err
	}
	r.natsConn = connect
	subscribe, err := r.natsConn.QueueSubscribe(subject, subject, func(msg *nats.Msg) {
		r.lastHandelTs = time.Now()
		respMsg := nats.NewMsg(msg.Subject)
		rspData, err2 := r.call(msg)
		if err2 != nil {
			respMsg.Header.Set("code", "-1")
			respMsg.Header.Set("msg", err2.Error())
		} else {
			respMsg.Data = rspData
		}
		respMsg.Header.Set("code", "0")
		err2 = msg.RespondMsg(respMsg)
		if err2 != nil {
			logrus.Errorf("connectNats RespondMsg err:%s", err2)
		}
	})
	r.natsSubscribe = subscribe
	if err != nil {
		logrus.Errorf("connectNats QueueSubscribe subject:%s uuid:%s err:%s", subject, r.uuid, err.Error())
		//logrus.Errorf("connectNats QueueSubscribe err:%s", err.Error())
		return err
	}

	msg := nats.NewMsg(r.uuid)
	msg.Header.Set("code", "0")
	respMsg, err := r.natsConn.RequestMsg(msg, time.Second*5)
	if err != nil {
		logrus.Errorf("connectNats RequestMsg subject:%s uuid:%s err:%s", subject, r.uuid, err.Error())
		return err
	}
	if respMsg.Header.Get("code") == "0" {
		logrus.Infof("connectNats subject:%s connect success cost:%s", subject, time.Now().Sub(now).String())
	} else {
		errMsg := respMsg.Header.Get("msg")
		logrus.Infof("connectNats subject:%s connect fail err:%s cost:%s", subject, errMsg, time.Now().Sub(now).String())
	}

	return nil
}

func (r *Runner) close() error {
	if r.isClosed {
		return nil
	}
	r.isClosed = true
	now := time.Now()
	subject := "close.runner"
	newMsg := nats.NewMsg(subject)
	newMsg.Header.Set("version", r.detail.Version)
	newMsg.Header.Set("user", r.detail.User)
	newMsg.Header.Set("name", r.detail.Name)
	newMsg.Data = []byte(r.uuid)

	msg, err := r.natsConn.RequestMsg(newMsg, time.Second*5)
	if err != nil {
		logrus.Infof("Runner close RequestMsg subject:%s uid:%s err:%s cost:%s\n", subject, r.uuid, err.Error(), time.Now().Sub(now).String())
		return err
	}
	if msg.Header.Get("code") == "0" {
		logrus.Infof("Runner close subject:%s uid:%s success cost:%s\n", subject, r.uuid, time.Now().Sub(now).String())
		return nil
	}
	logrus.Errorf("Runner close subject:%s uid:%s err:%s cost:%s\n", subject, r.uuid, msg.Header.Get("msg"), time.Now().Sub(now).String())

	return fmt.Errorf(msg.Header.Get("msg"))
}
