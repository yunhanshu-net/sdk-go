package runner

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"github.com/yunhanshu-net/sdk-go/model/request"
	v2 "github.com/yunhanshu-net/sdk-go/model/response/v2"
	"net"
	"os"
	"time"
)

//func (r *Runner) connect() error {
//
//	msg := nats.NewMsg(fmt.Sprintf("runcher.%s.%s.%s.connect",
//		r.detail.User, r.detail.Name, r.detail.Version))
//	msg.Header.Set("connect", "req")
//	msg.Header.Set("uuid", r.uuid)
//	msg.Header.Set("subject", r.detail.GetSubject())
//
//	c, err := net.Dial("unix", r.detail.GetAddr())
//	if err != nil {
//		return err
//	}
//	r.netConn = NewMessageConn(c)
//	buffer := make([]byte, 1024)
//	for {
//
//		data, err := r.netConn.ReadAll()
//		if err != nil {
//			fmt.Println("\nServer connection closed")
//		}
//		r.handelData(data)
//		fmt.Printf("\nServer response: %s", string(buffer[:n]))
//		fmt.Print("Enter message: ")
//	}
//
//	connect, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		logrus.Infof("connect:uuid failed: %s", r.uuid)
//		return err
//	}
//	logrus.Info("connect:uuid success: " + r.uuid)
//	r.conn = connect
//
//	group := fmt.Sprintf("%s.%s.%s", r.detail.User, r.detail.Name, r.detail.Version)
//
//	sub, err := connect.QueueSubscribe("runner.>", group, func(msg *nats.Msg) {
//		r.lastHandelTs = time.Now()
//		var reqMsg request.RunnerRequest
//		err1 := sonic.Unmarshal(msg.Data, &reqMsg)
//		if err1 != nil {
//			panic(err1)
//		}
//		//ctx := &Context{req: &reqMsg, Request: reqMsg.Request, ResponseData: &response.ResponseData{}}
//		httpContext := &HttpContext{
//			Request:  reqMsg.Request,
//			runner:   reqMsg.Runner,
//			Response: &v2.ResponseData{},
//		}
//		err = r.runRequest(httpContext)
//		marshal, err1 := sonic.Marshal(httpContext.Response)
//		if err1 != nil {
//			panic(err)
//		}
//		newMsg := nats.NewMsg(msg.Subject)
//		newMsg.Data = marshal
//		err1 = msg.RespondMsg(newMsg)
//		if err1 != nil {
//			panic(err1)
//		}
//	})
//	if err != nil {
//		msg.Header.Set("code", "-1")
//		msg.Header.Set("msg", err.Error())
//		_, _ = r.conn.RequestMsg(msg, time.Second*2)
//		panic(err)
//
//	}
//	r.sub = sub
//
//	_, err = r.conn.RequestMsg(msg, time.Second*2)
//	if err != nil {
//		logrus.Infof("connect RequestMsg Ping uuid:%s err:%s", r.uuid, err)
//		return err
//	}
//
//	logrus.Infof("connect done")
//	return nil
//}

type Rpc struct {
	r *Runner
}

func (r *Runner) GetUnixPath() string {
	return r.detail.GetUnixPath()
}
func (r *Runner) connectRpc() error {
	unixPath := r.GetUnixPath()
	os.Remove(unixPath)
	logrus.Infof("connectRpc" + unixPath)
	s := server.NewServer()
	rpc := &Rpc{r: r}
	err := s.Register(rpc, "")
	if err != nil {
		logrus.Errorf("connectRpc err:" + err.Error())
		return err
	}
	logrus.Infof("connectRpc success" + unixPath)
	r.rpcSrv = s
	fmt.Println("<connect-ok></connect-ok>")
	err = s.Serve("unix", unixPath)
	if err != nil {
		logrus.Error("connectRpc err:%s", err)
	}
	return nil
}

func (r *Rpc) Ping(ctx context.Context, req *request.Ping, response *request.Ping) error {
	clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	r.r.rpcConn = clientConn
	return nil
}

func (r *Rpc) Call(ctx context.Context, req *request.RunnerRequest, response *v2.ResponseData) error {
	r.r.lastHandelTs = time.Now()
	var err error
	httpContext := &HttpContext{
		Request:  req.Request,
		runner:   req.Runner,
		Response: response}
	//defer func() {
	//	logrus.Infof("call err:%+v req:%+v runner:%+v rsp:%+v",
	//		err, req.Request, req.Runner, httpContext.Response)
	//}()
	err = r.r.runRequest(httpContext)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rpc) Close(ctx context.Context, req *request.RunnerRequest, response *v2.ResponseData) error {
	logrus.Infof("call close:%s", r.r.GetUnixPath())
	response.StatusCode = 200
	response.Msg = "ok"
	r.r.down <- struct{}{}
	return nil
}

func (r *Runner) close() error {
	err := r.rpcSrv.SendMessage(r.rpcConn, "", "", nil, []byte(r.detail.GetUnixPath()))
	if err != nil {
		return err
	}
	defer r.rpcSrv.Close()

	//msg := nats.NewMsg(fmt.Sprintf("runcher.%s.%s.%s.close",
	//	r.detail.User, r.detail.Name, r.detail.Version))
	//msg.Header.Set("close", "req")
	//msg.Header.Set("subject", r.detail.GetSubject())
	//msg.Header.Set("uuid", r.uuid)
	//_, err := r.conn.RequestMsg(msg, time.Second*2)
	//if err != nil {
	//	return err
	//}
	//r.sub.Unsubscribe()
	//r.conn.Close()
	return nil
}
