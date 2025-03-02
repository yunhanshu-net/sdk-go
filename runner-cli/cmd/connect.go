package cmd

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"github.com/yunhanshu-net/sdk-go/runner"
	"strconv"
	"time"
)

func Connect(args []string) {
	fmt.Println("Connect:", args)
	//return
	runnerName := args[0] //runner的绝对路径
	route := args[1]      //
	reqFile := args[2]    //路径
	count := 1

	var req runner.Request
	err := jsonx.UnmarshalFromFile(reqFile, &req)
	if err != nil {
		panic(err)
	}
	argsMap := make(map[string]bool)
	for _, arg := range args {
		argsMap[arg] = true
	}
	fmt.Printf("count: %v\n", args[3])
	if len(args) > 3 {
		i, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			panic(err)
		}
		count = int(i)
	}
	now := time.Now()
	run(runnerName, route, reqFile, argsMap["-no-out"])
	fmt.Printf("run task count:%v total cost:%s\n", count, time.Since(now))
}
