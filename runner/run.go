package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"os"
	"runtime"
)

func (r *Runner) Run() {
	runtime.GOMAXPROCS(1)
	command := os.Args[1]
	jsonFileName := os.Args[2]
	fmt.Println(command)
	fmt.Println(jsonFileName)
	//workPath := os.Args[3]
	dir, err := os.Getwd()
	fmt.Println("当前目录。", dir)
	fmt.Println(err)

	m := runtime.MemStats{}
	//todo
	runtime.ReadMemStats(&m)

	//todo
	reqContext := &Context{
		Request: &request.Request{},
	}
	fmt.Println(reqContext)

	//todo handel ctx
	//context := &Context{Request: jsonFileName, IsDebug: r.IsDebug, Cmd: command, runtimeInfo: &RuntimeInfo{
	//	start: m,
	//}}

	//Handel(context, r)
	//if r.About {
	//	return
	//}

}
