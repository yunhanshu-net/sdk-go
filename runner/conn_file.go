package runner

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"sync"
)

type connResponse struct {
	Response *response.Response     `json:"response"`
	MetaData map[string]interface{} `json:"meta_data"` //调用耗时，内存占用，日志信息等等
}

type fileConn struct {
	args      []string
	info      *Info
	reqWg     *sync.WaitGroup
	nats      *nats.Conn
	sub       *nats.Subscription
	requestCh chan *request.Request
}

func (f *fileConn) Connect() error {
	var req request.Request
	err := jsonx.UnmarshalFromFile(f.args[2], &req)
	if err != nil {
		fmt.Println("jsonx.UnmarshalFromFile(jsonFileName, &req) err:" + err.Error())
		return err
	}
	return nil
}

func (f *fileConn) Response(resp *response.Response) error {
	rsp, err := json.Marshal(&connResponse{resp, nil})
	if err != nil {
		return err
	}
	s := string(rsp)
	fmt.Println("<Response>" + s + "</Response>")
	return nil
}

func (f *fileConn) Close() error {
	return nil
}
