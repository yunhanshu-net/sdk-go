package syncx

import "sync"

// ConcurrencyControl 并发控制函数
func ConcurrencyControl(tasks []func(), limit int) {
	// 创建一个带缓冲的通道，用于控制并发数
	semaphore := make(chan struct{}, limit)
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1) // 增加WaitGroup计数
		go func(task func()) {
			defer wg.Done() // 确保goroutine结束时减少计数

			semaphore <- struct{}{} // 占用一个并发槽位
			task()                  // 执行任务
			<-semaphore             // 释放一个并发槽位
		}(task)
	}

	wg.Wait() // 等待所有任务完成
}
