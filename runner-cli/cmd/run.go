package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func Run(args []string) {
	fmt.Println("run:", args)
	//return
	runner := args[0]  //runner的绝对路径
	route := args[1]   //
	reqFile := args[2] //路径
	count := 1

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
	isSync := argsMap["-sync"]
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for range count {
		if isSync {
			go func() {
				run(runner, route, reqFile, argsMap["-no-out"])
				defer wg.Done()
			}()
		} else {
			run(runner, route, reqFile, argsMap["-no-out"])
		}
	}
	if isSync {
		wg.Wait()
	}
	fmt.Printf("run task count:%v total cost:%s\n", count, time.Since(now))

}

func run(runner, route, reqFile string, noOut bool) {
	cmd := exec.Command(runner, route, reqFile)
	// 捕获命令的输出
	var out bytes.Buffer
	cmd.Stdout = &out

	now := time.Now()
	// 执行命令
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Route execution %s %s %s failed: %v\n", runner, route, reqFile, err)
		// 打印输出结果
		fmt.Println("Route Output:")
		fmt.Println(out.String())
		return
	}

	if !noOut {
		// 打印输出结果
		fmt.Printf("Route Output: cost:%s\n", time.Since(now))
		fmt.Println(out.String())
	}

}
