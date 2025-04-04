package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/syncx"
	"log"
	"time"

	"github.com/smallnest/rpcx/client"
)

var (
	addr = "/Users/yy/Desktop/code/github.com/sdk-go/test/rpcxs/srv/test.sock"
)

func main() {
	flag.Parse()

	d, _ := client.NewPeer2PeerDiscovery("unix@"+addr, "")
	//client.DefaultOption.SerializeType = protocol.JSON
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	Test(xclient)

}
func Test(xClient client.XClient) {
	now := time.Now()
	n := 100000
	tasks := make([]func(), 0, n)
	for i := 0; i < n; i++ {
		tasks = append(tasks, func() {
			args := &request.RunnerRequest{Request: &request.Request{
				Method: "GET",
				Route:  "api",
				Body: map[string]interface{}{
					"test":  1,
					"hello": "world",
				},
			}, Runner: &model.Runner{
				Name:    "openapi",
				User:    "tencent",
				Version: "v1",
			}}
			reply := &response.RunnerResponse{}
			err := xClient.Call(context.Background(), "Mul", args, reply)
			if err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			//fmt.Printf("%+vbody:%+v\n", *reply, reply.Response.Body)
		})
	}
	syncx.ConcurrencyControl(tasks, 100)
	fmt.Println(time.Since(now).String())
}
