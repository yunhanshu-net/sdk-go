package main

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/syncx"
	"log"
	"time"
)

func main() {
	addr := "/Users/yy/Desktop/code/github.com/sdk-go/soft/beiluo/debug/beiluo_debug_v1.sock"
	d, _ := client.NewPeer2PeerDiscovery("unix@"+addr, "")
	//client.DefaultOption.SerializeType = protocol.JSON
	xclient := client.NewXClient("Runner", client.Failtry, client.RandomSelect, d, client.DefaultOption)
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
				Route:  "/hello",
				Body: map[string]interface{}{
					"test":  1,
					"hello": "world",
				},
			}, Runner: &model.Runner{
				Name:    "debug",
				User:    "beiluo",
				Version: "v1",
			}}
			reply := &response.Data{}
			err := xClient.Call(context.Background(), "call", args, reply)
			if err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			//fmt.Printf("%+vbody:%+v\n", *reply, reply.Body)
		})
	}
	syncx.ConcurrencyControl(tasks, 100)
	fmt.Println(time.Since(now).String())
}
