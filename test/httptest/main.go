package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func httpRequest(url string, wg *sync.WaitGroup) {
	//defer wg.Done() // 在函数结束时调用 Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(all))
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	fmt.Printf("请求 %s 返回状态码: %d\n", url, resp.StatusCode)
}

func loadTest(url string, numRequests int, limit int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, limit)
	start := time.Now() // 记录开始时间
	wg.Add(numRequests) // 增加 WaitGroup 计数
	for i := 0; i < numRequests; i++ {
		// 增加 WaitGroup 计数
		semaphore <- struct{}{} // 获取信号量
		go func() {
			defer func() {

				<-semaphore // 释放信号量
				wg.Done()
			}()
			httpRequest(url, &wg) // 发送请求
		}()
		//go httpRequest(url, &wg) // 启动 goroutine 发送请求
	}

	wg.Wait()                     // 等待所有请求完成
	duration := time.Since(start) // 计算总耗时
	fmt.Printf("总请求数: %d, 耗时: %v\n", numRequests, duration)
}

func main() {

	url := "http://127.0.0.1:8888/runner/beiluo/apphub/hello"
	numRequests := 3000

	loadTest(url, numRequests, 10)
}
