package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

var Conn *nats.Conn
var Sub *nats.Subscription
var ReceiveSub *nats.Subscription

func InitConn() {
	connect, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	Conn = connect

	subscribe, err := connect.Subscribe("runcher.>", func(msg *nats.Msg) {
		marshal, err1 := json.Marshal(msg.Header)
		if err1 != nil {
			fmt.Println(err1)
		}
		newMsg := nats.NewMsg(msg.Subject)
		newMsg.Data = []byte("res ok")
		//err1 = connect.Publish(msg.Reply, []byte("res ok"))
		err1 = msg.RespondMsg(newMsg)
		if err1 != nil {
			panic(err1)
		}
		fmt.Printf("receive <- subject:%s\ndata:\n%s\nheaders:%s\n", msg.Subject, string(msg.Data), string(marshal))
	})

	//ReceiveSub, err = connect.Subscribe("runner.>", func(msg *nats.Msg) {
	//	//marshal, err1 := json.Marshal(msg.Header)
	//	//if err1 != nil {
	//	//	fmt.Println(err1)
	//	//}
	//	//newMsg := nats.NewMsg(msg.Subject)
	//	//newMsg.Data = []byte("res ok")
	//	////err1 = connect.Publish(msg.Reply, []byte("res ok"))
	//	////err1 = msg.RespondMsg(newMsg)
	//	//if err1 != nil {
	//	//	panic(err1)
	//	//}
	//	fmt.Printf("receive <- subject:%s\ndata:\n%s\n", msg.Subject, string(msg.Data))
	//})

	if err != nil {
		panic(err)
	}
	Sub = subscribe
	Conn = connect
	//defer subscribe.Unsubscribe()

}
