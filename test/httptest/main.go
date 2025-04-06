package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type T struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func todo(interface{}) {

}
func httpRequest(url string, wg *sync.WaitGroup) error {
	//defer wg.Done() // 在函数结束时调用 Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("code:%v", resp.StatusCode)
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	todo(all)
	var t T
	err = json.Unmarshal(all, &t)
	if err != nil {
		return err
	}
	if t.Code != 0 {
		return fmt.Errorf(t.Msg)
	}

	s := string(all)
	fmt.Println(s)

	//fmt.Printf("请求 %s 返回状态码: %d\n", url, resp.StatusCode)
	return nil
}

func loadTest(url string, numRequests int, limit int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, limit)
	res := make(chan error, numRequests)
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
			err := httpRequest(url, &wg) // 发送请求
			if err != nil {
				res <- err
			}
		}()
		//go httpRequest(url, &wg) // 启动 goroutine 发送请求
	}

	wg.Wait() // 等待所有请求完成
	if len(res) == 0 {
		fmt.Println("无错误")
	} else {
		for err := range res {
			fmt.Println(err.Error())
		}
	}
	duration := time.Since(start) // 计算总耗时
	fmt.Printf("总请求数: %d, 耗时: %v\n", numRequests, duration)
}

func main() {

	url := "http://127.0.0.1:9999/api/runner/beiluo/debug/hello"
	numRequests := 1000

	loadTest(url, numRequests, 50)
}
