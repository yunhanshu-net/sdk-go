package runner

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
	"os"
	"time"
)

func init() {
	logger.Setup()
}

func (r *Runner) listen() {
	ticker := time.NewTicker(time.Second * 1)
	logrus.Infof("listen uuid:%s\n", r.uuid)
	defer ticker.Stop()
	for {
		select {
		//case <-r.down:
		//	r.close()
		//	logrus.Infof("%s runcher发起关闭请求，已经关闭连接", r.GetUnixPath())
		//	return
		case <-ticker.C:
			if r.idle > 0 {
				ts := time.Now().Unix()
				d := ts - r.lastHandelTs.Unix()
				if (ts - r.lastHandelTs.Unix()) > r.idle { //超过指定空闲时间的话需要释放进程
					logrus.Infof(" %v没有处理消息，runner 自动关闭连接 idle config：%v", d, r.idle)
					r.close()
					return
				}
			}

		}
	}
}

func (r *Runner) Debug(user, runner, version string, idle int64, uuid string) error {
	r.uuid = uuid
	r.isDebug = true
	r.detail = &model.Runner{}
	r.detail.Name = runner
	r.detail.User = user
	r.detail.Version = version
	r.idle = idle
	go func() {
		err := r.connectNats()
		if err != nil {
			logrus.Error(err)
		}
	}()
	r.listen()
	return nil
}

func (r *Runner) Run() error {
	err := r.init(os.Args)
	if err != nil {
		return err
	}
	return nil
}
