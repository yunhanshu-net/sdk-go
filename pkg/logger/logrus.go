// Package logger ...
package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"time"
)

// Option ...
type Option struct {
}

// Setup ...
func Setup(option ...*Option) {
	mw := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    50,             // 文件最大大小（MB）
		MaxBackups: 10,             // 保留旧文件的最大数量
		MaxAge:     30,             // 保留旧文件的最大天数
		Compress:   false,          // 是否压缩旧文件
	}
	logrus.SetOutput(mw)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.DateTime})
}
