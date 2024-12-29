package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"os"
)

const (
	commandConnect = "_connect_"
)

func (r *Runner) start() {
	command := os.Args[1]
	jsonFileName := os.Args[2]
	r.init()

	switch command {
	case commandConnect: //建立长连接
		nt := &natsConn{}
		err := nt.Connect(r)
		if err != nil {
			return
		}
		//connect, err := nats.Connect(nats.DefaultURL)
		//if err != nil {
		//	panic(err)
		//}
		//r.nats = connect
		////todo 连接nats
		//err = r.connect()
		//if err != nil {
		//	return
		//}
	default:
		//默认即时调用模式，调用后立刻结束程序
		var req request.Request
		err := jsonx.UnmarshalFromFile(jsonFileName, &req)
		if err != nil {
			fmt.Println("jsonx.UnmarshalFromFile(jsonFileName, &req) err:" + err.Error())
			return
		}
		f := &fileIo{}
		r.handel(f)
		r.exit = make(<-chan struct{})
	}

}

func (r *Runner) Run() {

	r.start()

	select {
	case <-r.exit:
		//todo 释放资源
		if r.isKeepAlive {
			r.wg.Wait()
			r.Close()
		}
		return
	}
	//单核心执行，防止程序把资源占用尽

	//command := os.Args[1]
	//jsonFileName := os.Args[2]
	//fmt.Println(command)
	//fmt.Println(jsonFileName)
	////workPath := os.Args[3]
	//dir, err := os.Getwd()
	//fmt.Println("当前目录。", dir)
	//fmt.Println(err)

	//m := runtime.MemStats{}
	////todo
	//runtime.ReadMemStats(&m)
	//
	////todo
	//reqContext := &Context{
	//	Request: &request.Request{},
	//}
	//fmt.Println(reqContext)

	//for {
	//	select {
	//	case ctx := <-r.contextChan:
	//		fmt.Println(ctx)
	//	case <-r.exit:
	//		return
	//
	//	}
	//}

	//todo handel ctx
	//context := &Context{Request: jsonFileName, IsDebug: r.IsDebug, Cmd: command, runtimeInfo: &RuntimeInfo{
	//	start: m,
	//}}

	//Handel(context, r)
	//if r.About {
	//	return
	//}
}
