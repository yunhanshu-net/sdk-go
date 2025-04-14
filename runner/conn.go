package runner

import (
	"context"
	"fmt"
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

func (r *Runner) close() error {
	err := r.rpcSrv.SendMessage(r.rpcConn, "", "", map[string]string{"type": "close"}, []byte("close"))
	if err != nil {
		return err
	}
	defer r.rpcSrv.Close()
	return nil
}
