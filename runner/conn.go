package runner

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"net"
	"os"
	"time"
)

type Rpc struct {
	r *Runner
}

func (r *Runner) GetUnixPath() string {
	return r.detail.GetUnixPath()
}
func (r *Runner) connectRpc() error {
	unixPath := r.GetUnixPath()
	os.Remove(unixPath)
	s := server.NewServer()
	rpc := &Rpc{r: r}
	err := s.Register(rpc, "")
	if err != nil {
		logrus.Errorf("connectRpc err:" + err.Error())
		return err
	}
	r.rpcSrv = s
	fmt.Println("<connect-ok></connect-ok>")
	err = s.Serve("unix", unixPath)
	if err != nil {
		logrus.Errorf("connectRpc Serve err:%s", err)
	}
	return nil
}

func (r *Runner) call(msg *nats.Msg) ([]byte, error) {

	//router, exist := r.getRouter(req.Request.Route, req.Request.Method)
	//if !exist {
	//	logrus.Errorf("call getRouter 404,%s %s", req.Request.Method, req.Request.Route)
	//	return nil, fmt.Errorf("call getRouter 404,%s %s", req.Request.Method, req.Request.Route)
	//}
	data := msg.Data
	var req request.RunnerRequest
	err1 := sonic.Unmarshal(data, &req)
	if err1 != nil {
		logrus.Errorf("call  sonic.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call  sonic.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
	}

	runRequest, err1 := r.runRequest(context.Background(), req.Request)
	if err1 != nil {
		logrus.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
	}
	marshal, err1 := sonic.Marshal(runRequest)
	if err1 != nil {
		logrus.Errorf("call sonic.Marshal err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call sonic.Marshal err,req:%+v err:%s", req, err1.Error())
	}

	return marshal, nil
}

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

func (r *Rpc) Ping(ctx context.Context, req *request.Ping, response *request.Ping) error {
	clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	r.r.rpcConn = clientConn
	return nil
}

func (r *Rpc) Call(ctx context.Context, req *request.RunnerRequest, resp *response.Data) error {
	r.r.lastHandelTs = time.Now()
	var err error
	//httpContext := &HttpContext{Request: req.Request, runner: req.Runner, Response: resp}
	rsp, err := r.r.runRequest(ctx, req.Request)
	if err != nil {
		return err
	}
	*resp = *rsp
	//todo 判断是否需要reset body
	return nil
}

func (r *Rpc) Close(ctx context.Context, req *request.RunnerRequest, response *response.Data) error {
	logrus.Infof("call close:%s", r.r.GetUnixPath())
	response.StatusCode = 200
	response.Msg = "ok"
	r.r.down <- struct{}{}
	return nil
}

//
//func (r *Runner) close() error {
//	err := r.rpcSrv.SendMessage(r.rpcConn, "", "", map[string]string{"type": "close"}, []byte("close"))
//	if err != nil {
//		return err
//	}
//	defer r.rpcSrv.Close()
//	return nil
//}

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
