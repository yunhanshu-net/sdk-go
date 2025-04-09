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
	return nil
}
