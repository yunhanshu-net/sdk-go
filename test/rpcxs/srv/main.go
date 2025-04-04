package main

import (
	"context"
	"flag"
	"github.com/smallnest/rpcx/server"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"os"
)

type Arith struct{}

// the second parameter is not a pointer
func (t *Arith) Mul(ctx context.Context, req request.RunnerRequest, resp *response.RunnerResponse) error {
	//fmt.Printf("req.Request:%+v\n", req.Request)
	//fmt.Printf("req.Runner:%+v\n", req.Runner)
	//resp = new(response.RunnerResponse)
	resp.MetaData = make(map[string]interface{})
	resp.Response = new(response.Response)
	resp.Response.MetaData = make(map[string]interface{})
	resp.MetaData["test"] = 1
	resp.Response.MetaData["cost"] = 1
	resp.Response.MetaData["success"] = true
	resp.Response.Body = map[string]interface{}{"code": 0, "msg": "ok"}
	return nil
}

func main() {
	flag.Parse()
	sock := "/Users/yy/Desktop/code/github.com/sdk-go/test/rpcxs/srv/test.sock"
	os.Remove(sock)
	s := server.NewServer()
	defer s.Close()
	err := s.Register(new(Arith), "")
	if err != nil {
		panic(err)
	}
	err = s.Serve("unix", sock)
	if err != nil {
		panic(err)
	}
}
