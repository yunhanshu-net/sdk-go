package runner

// 先建立三次握手 确保连接上nats，双方能够正常通信
//func (r *Runner) connect() error {
//	//调度引擎发起第一次握手
//	connectKey := fmt.Sprintf("runner.connect.%s.%s", r.info.User, r.info.Runner)
//	connectMsg := nats.NewMsg(connectKey)
//	handelKey := fmt.Sprintf("runner.req.%s.%s", r.info.User, r.info.Runner)
//	sub, err := r.nats.Subscribe(handelKey, r.onMsg)
//	if err != nil { //连接失败响应
//		connectMsg.Header.Set("status", "-1")
//		connectMsg.Header.Set("msg", "connect_subscribe_failed")
//		connectMsg.Header.Set("error", err.Error())
//		_, err := r.nats.RequestMsg(connectMsg, time.Second*2) //连接失败
//		if err != nil {
//			return err
//		}
//		return err
//	}
//	r.sub = sub
//
//	//连接成功响应
//	connectMsg.Header.Set("status", "0")
//	connectMsg.Header.Set("msg", "success")
//	//发送连接报文响应给调度引擎，第二次握手
//	rspConnectMsg, err := r.nats.RequestMsg(connectMsg, time.Second*3)
//	if err != nil {
//		return err
//	}
//
//	//确认调度引擎的回复信息，第三次握手
//	if rspConnectMsg.Header.Get("status") != "0" { //说明失败
//		return fmt.Errorf(rspConnectMsg.Header.Get("msg"))
//	}
//	//建立三次握手，连接成功，此时已经可以正常消费到消息了
//	return nil
//}
