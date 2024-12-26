package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"runtime"
)

const (
	commandConnect = "connect"
)

func (r *Runner) ReadRequest() {

}

func (r *Runner) init() {
	runtime.GOMAXPROCS(1)
}

func (r *Runner) start() {
	command := os.Args[1]
	//jsonFileName := os.Args[2]

	switch command {
	case commandConnect: //建立长连接
		connect, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			panic(err)
		}
		r.nats = connect
		//todo 连接nats
		err = r.connect()
		if err != nil {
			return
		}
	default:
		//默认即时调用模式，调用后立刻结束程序

	}

}

func (r *Runner) Run() {

	//单核心执行，防止程序把资源占用尽

	command := os.Args[1]
	jsonFileName := os.Args[2]
	fmt.Println(command)
	fmt.Println(jsonFileName)
	//workPath := os.Args[3]
	dir, err := os.Getwd()
	fmt.Println("当前目录。", dir)
	fmt.Println(err)

	//m := runtime.MemStats{}
	////todo
	//runtime.ReadMemStats(&m)
	//
	////todo
	//reqContext := &Context{
	//	Request: &request.Request{},
	//}
	//fmt.Println(reqContext)

	for {
		select {
		case ctx := <-r.contextChan:
			fmt.Println(ctx)
		case <-r.exit:
			return

		}
	}

	//todo handel ctx
	//context := &Context{Request: jsonFileName, IsDebug: r.IsDebug, Cmd: command, runtimeInfo: &RuntimeInfo{
	//	start: m,
	//}}

	//Handel(context, r)
	//if r.About {
	//	return
	//}
}
