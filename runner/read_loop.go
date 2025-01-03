package runner

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Runner) readLoop() {
	r.lastHandelTs = time.Now().Unix()
	//检查连接是否活跃，不活跃的情况下关闭连接
	tk := time.NewTicker(time.Second * 5)
	for {
		select {
		case req, ok := <-r.requestCh:
			r.lastHandelTs = time.Now().Unix()
			//fmt.Println("readLoop", req)
			req.TraceID = uuid.New().String()
			go func() { //每读到一个请求就异步处理
				ctx, err := r.handelRequest(req)
				if err != nil {
					fmt.Println("readLoop", err)
				}
				r.conn.Response(ctx.Response)
			}()
			if !ok {
				return
			}
		case <-tk.C:
			if time.Now().Unix()-r.lastHandelTs > 5 {
				r.Close() //先发送关闭信号
				close(r.requestCh)
				//todo 这里应该先close然后再sleep一下，确保消息消费完毕
				r.exit()
			}
		default:
		}
	}
}
