package main

import (
	"bytes"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/unix"
	"net"
	"os"
)

const socketPath = "/tmp/echo.sock"

func main() {
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	mc := unix.NewMessageConn(conn)
	defer mc.Close()

	for {
		// 读取客户端请求
		id, data, err := mc.readWithUUID()
		if err != nil {
			return
		}

		// 处理请求（示例：转大写）
		go func() {
			respData := bytes.ToUpper(data)
			mc.writeWithUUID(id, respData)
		}()
	}
}

// 客户端收到服务端的请求
func clientOnRequest() {
	conn.OnRequest(func(msg) {
		key := msg.GetHeaders("key")
		fmt.Println(key) //server set key
		data := msg.GetData()
		rsp := newMsg(string(data) + "client resp")
		msg.Respone(rsp)
	})
}

// 服务端请求客户端
func server2client() {
	reqMsg可以设置header
	reqMsg = newMsg("hello")
	reqMsg.SetHeader("key", "server set key")
	resp, err := conn.Request(reqMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Data)) //hello client resp
}

// 服务端请求客户端
func client2server() {
	reqMsg可以设置header
	reqMsg = newMsg("hello 我要关闭连接")
	reqMsg.SetHeader("key", "client set key")
	resp, err := conn.Request(reqMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Data)) //hello 我要关闭连接server resp
}

// 客户端收到服务端的请求
func serverOnRequest() {
	conn.OnRequest(func(msg) {
		key := msg.GetHeaders("key")
		fmt.Println(key) //client set key
		data := msg.GetData()
		rsp := newMsg(string(data) + "server resp")
		msg.Respone(rsp)
	})
}
