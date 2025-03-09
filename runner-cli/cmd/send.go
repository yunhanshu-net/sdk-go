package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"strconv"
	"sync"
	"time"
)

func Send(args []string) {

	fmt.Println("run:", args)
	//return
	runnerName := args[0] //runner的绝对路径
	//route := args[1]   //
	reqFile := args[1] //路径
	count := 1

	argsMap := make(map[string]bool)
	for _, arg := range args {
		argsMap[arg] = true
	}
	fmt.Printf("count: %v\n", args[2])
	if len(args) > 2 {
		i, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			panic(err)
		}
		count = int(i)
	}
	var req request.RunnerRequest
	err := jsonx.UnmarshalFromFile(reqFile, &req)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	isSync := argsMap["-sync"]
	wg := &sync.WaitGroup{}
	wg.Add(count)
	c := count
	lk := &sync.RWMutex{}
	for range count {
		if isSync {
			go func() {
				send(runnerName, &req, argsMap["-no-out"])
				lk.Lock()
				c--
				fmt.Printf("%v\n", c)
				lk.Unlock()
				defer wg.Done()
			}()
		} else {
			send(runnerName, &req, argsMap["-no-out"])
		}
	}
	if isSync {
		wg.Wait()
	}
	fmt.Printf("send task count:%v total cost:%s\n", count, time.Since(now))
}

func send(runner string, req *request.RunnerRequest, noOut bool) {
	subject := fmt.Sprintf("upstream.%s.%s.%s.run", req.Runner.User, req.Runner.Name, req.Runner.Version)
	//fmt.Printf("subject:%s\n", subject)
	msg := nats.NewMsg(subject)
	marshal, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	msg.Data = marshal
	msg.Header.Set("user", req.Runner.User)
	msg.Header.Set("runner", req.Runner.Name)
	msg.Header.Set("version", req.Runner.Version)
	msg.Header.Set("method", req.Request.Method)
	msg.Header.Set("route", req.Request.Route)
	requestMsg, err := Conn.RequestMsg(msg, time.Second*50)
	if err != nil {
		panic(err)
	}
	if !noOut {
		fmt.Printf("send to -> subject:%s data:%s\n", subject, string(msg.Data))
		fmt.Printf("recv to -> subject:%s data:%s\n", requestMsg.Subject, string(requestMsg.Data))
	}

}
