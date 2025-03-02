package runner

import (
	"github.com/yunhanshu-net/sdk-go/model/response"
)

//func NewConn(r *Runner) Conn {
//	if r.isKeepAlive {
//		return &natsConn{info: r.info, reqWg: r.wg, requestCh: r.requestCh}
//	}
//	return &fileConn{args: r.args, info: r.info, reqWg: r.wg, requestCh: r.requestCh}
//}

type Conn interface {
	Connect() error
	Response(resp *response.Response) error
	Close() error
}
