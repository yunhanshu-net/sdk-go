package main

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/runner-cli/cmd"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	// 创建 readline 实例
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",     // 提示符
		HistoryFile:     ".history", // 历史记录文件（可选）
		InterruptPrompt: "^C",       // 按 Ctrl+C 时的提示
		EOFPrompt:       "exit",     // 按 Ctrl+D 时的提示
	})
	if err != nil {
		log.Fatalf("Failed to initialize readline: %v", err)
	}
	defer rl.Close()

	fmt.Println("Welcome to the Interactive Terminal with History Support!")
	fmt.Println("Type 'exit' or press Ctrl+D to quit.")

	cmd.InitConn()
	defer cmd.Conn.Close()
	defer cmd.Sub.Unsubscribe()
	for {
		// 获取用户输入
		line, err := rl.Readline()
		if err != nil { // 检查是否退出
			if err == readline.ErrInterrupt { // 用户按 Ctrl+C
				fmt.Println("^C")
				continue
			}
			log.Fatalf("Error reading input: %v", err)
		}

		// 去除首尾空白
		input := strings.TrimSpace(line)

		// 检查是否退出
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// 调用自定义处理函数
		handleInput(input)
	}
}

// handleInput 是你自定义的处理逻辑
func handleInput(input string) {
	split := strings.Split(input, ">")
	for i, s := range split {
		split[i] = strings.Trim(s, " ")
	}
	cmd.InitArgs(split)
}
