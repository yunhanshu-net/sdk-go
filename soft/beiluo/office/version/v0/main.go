package main

import (
	"github.com/yunhanshu-net/sdk-go/runner"
)

func InitPackages() {
}

func main() {
	InitPackages()
	err := runner.Run()
	if err != nil {
		panic(err)
	}
}
