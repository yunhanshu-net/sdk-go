package runner

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model"
	"os"
	"time"
)

func (r *Runner) listen() {
	ticker := time.NewTicker(time.Second * 1)
	logrus.Infof("listen uuid:%s\n", r.uuid)
	defer func() {
		ticker.Stop()
		// 使用统一的Shutdown函数而不是单独关闭资源
		Shutdown()
	}()

	for {
		select {
		case <-r.down:
			logrus.Infof("%s runcher发起关闭请求，关闭连接", r.uuid)
			return
		case <-ticker.C:
			if r.idle > 0 {
				ts := time.Now().Unix()
				d := ts - r.lastHandelTs.Unix()
				if (ts - r.lastHandelTs.Unix()) > r.idle { //超过指定空闲时间的话需要释放进程
					logrus.Infof(" %v没有处理消息，runner 自动关闭连接 idle config：%v", d, r.idle)
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

	// 创建一个channel来跟踪连接状态
	errChan := make(chan error, 1)
	go func() {
		err := r.connectNats(context.Background())
		if err != nil {
			errChan <- err
			logrus.Error(err)
		}
	}()

	// 等待连接成功或失败
	select {
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Second):
		// 继续执行
	}

	r.listen()
	return nil
}

func (r *Runner) Run() error {
	err := r.init(os.Args)

	// 注意：如果是连接模式，init方法内部会调用listen并阻塞，listen退出时会调用Shutdown
	// 如果是单次执行模式，runner.go中的run方法已经处理了Shutdown
	if err != nil {
		return err
	}

	// 不需要在这里调用Shutdown，因为：
	// 1. 连接模式下，listen()方法会在defer中调用Shutdown
	// 2. 单次执行模式下，run()方法已经调用了Shutdown
	//todo 那为啥我在这里统一释放？

	return nil
}
