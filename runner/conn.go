package runner

import "github.com/yunhanshu-net/sdk-go/model/request"

type Conn interface {
	Connect(runner *Runner) error
	RequestChan(requestCh chan<- *request.Request) (err error)
	Close() error
}
