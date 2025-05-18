package main

import (
	"github.com/yunhanshu-net/sdk-go/runner"
	"github.com/yunhanshu-net/sdk-go/soft/beiluo/te/version/v1/api/calc"
)

func init() {
	calc.Init()
}

func main() {
	runner.Run()
}
