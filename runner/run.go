package runner

import (
	"fmt"
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
		// 每次处理消息后重置定时器
		//idleTimer.Reset(timeout)
		select {
		case <-ticker.C:
			logrus.Infof("check uuid:%s\n", r.uuid)

			if r.idle > 0 {
				ts := time.Now().Unix()
				if (ts - r.lastHandelTs.Unix()) > r.idle { //超过指定空闲时间的话需要释放进程
					logrus.Infof("close uuid:%s\n", r.uuid)
					r.close()
					return
				}
			}

		}
	}
}

func (r *Runner) Debug(user, runner, version string, idle int64, uuid string) error {
	r.uuid = uuid
	r.detail = &model.Runner{}
	r.detail.Name = runner
	r.detail.User = user
	r.detail.Version = version
	r.idle = idle
	err := r.connect()
	if err != nil {
		return err
	}
	r.listen()
	return nil
}

func (r *Runner) Run() error {
	fmt.Println("handelFunctions:", r.handelFunctions)
	err := r.init(os.Args)
	if err != nil {
		return err
	}
	return nil
}
