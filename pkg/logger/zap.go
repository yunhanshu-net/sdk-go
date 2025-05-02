package logger

import (
	"context"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/constants"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 添加新的常量和变量
type contextKey string

const (
	// TraceIDKey 是存储在context中的traceid键
	TraceIDKey = contextKey("trace_id")
	// DefaultTraceID 当无法从context获取traceid时的默认值
	DefaultTraceID = "unknown"
)

var (
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
	baseDir string // 程序的基础目录
)

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 自定义调用者编码器
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 转换为相对路径
	path := getRelativePath(caller.File)

	// 获取函数名
	funcName := caller.Function
	if idx := strings.LastIndex(funcName, "."); idx > 0 {
		funcName = funcName[idx+1:]
	}

	// 格式化为"相对路径:行号 [函数名]"
	enc.AppendString(fmt.Sprintf("%s:%d [%s]", path, caller.Line, funcName))
}

// 获取相对路径，去除基础目录部分
func getRelativePath(path string) string {
	// 如果路径中包含basePath，只保留相对部分
	if baseDir != "" && strings.Contains(path, baseDir) {
		rel, err := filepath.Rel(baseDir, path)
		if err == nil {
			return rel
		}
	}

	// 如果无法获取相对路径，返回文件名和上级目录
	dir, file := filepath.Split(path)
	parent := filepath.Base(dir)
	return filepath.Join(parent, file)
}

func init() {
	// 获取程序基础目录
	_, file, _, ok := runtime.Caller(0)
	if ok {
		// 找到当前文件所在的目录
		dir := filepath.Dir(file)
		// 向上找两级目录（pkg/logger → runcher）
		baseDir = filepath.Dir(filepath.Dir(dir))
	}

	//isDev := os.Getenv("ENV") == "dev"

	// 创建基础配置
	config := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	if isDev {
		// 开发环境使用彩色日志级别
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置日志输出
	var core zapcore.Core
	if isDev {
		// 开发环境直接输出到控制台
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		core = zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zap.DebugLevel),
		)
	} else {
		// 生产环境输出到文件
		jsonEncoder := zapcore.NewJSONEncoder(config)
		writer := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    50,    // MB
			MaxBackups: 10,    // 保留旧文件的最大数量
			MaxAge:     30,    // 保留旧文件的最大天数
			Compress:   false, // 是否压缩旧文件
		}
		core = zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(writer),
			zap.NewAtomicLevelAt(zap.InfoLevel),
		)
	}

	// 创建logger实例
	logger = zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(1))
	sugar = logger.Sugar()
}

// extractTraceID 从context中提取trace_id
func extractTraceID(ctx context.Context) string {
	if ctx == nil {
		return DefaultTraceID
	}

	value := ctx.Value(constants.TraceID)
	// 尝试从context中获取trace_id
	traceID, ok := value.(string)

	if ok && traceID != "" {
		return traceID
	}

	return DefaultTraceID
}

// withTraceID 在日志字段中添加trace_id
func withTraceID(ctx context.Context, fields []zap.Field) []zap.Field {
	traceID := extractTraceID(ctx)
	if traceID != DefaultTraceID {
		fields = append(fields, zap.String(constants.TraceID, traceID))
	}
	return fields
}

// 以下是带有Context的日志方法

// DebugContext 输出带有trace_id的调试日志
func DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Debug(msg, withTraceID(ctx, fields)...)
}

// InfoContext 输出带有trace_id的信息日志
func InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(msg, withTraceID(ctx, fields)...)
}

// WarnContext 输出带有trace_id的警告日志
func WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Warn(msg, withTraceID(ctx, fields)...)
}

// ErrorContext 输出带有trace_id的错误日志
func ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(msg, withTraceID(ctx, fields)...)
}

// FatalContext 输出带有trace_id的致命错误日志
func FatalContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Fatal(msg, withTraceID(ctx, fields)...)
}

// DebugContextf 格式化输出带有trace_id的调试日志
func DebugContextf(ctx context.Context, template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		DebugContext(ctx, fmt.Sprintf(template, args...))
		return
	}

	traceID := extractTraceID(ctx)
	// 使用调用者的位置创建临时logger
	debugLogger := logger.WithOptions(zap.AddCallerSkip(0))
	if traceID != DefaultTraceID {
		// 输出日志，添加trace_id
		debugLogger.Debug(fmt.Sprintf(template, args...), zap.String("trace_id", traceID))
	} else {
		// 输出日志，不添加trace_id
		debugLogger.Debug(fmt.Sprintf(template, args...))
	}
}

// InfoContextf 格式化输出带有trace_id的信息日志
func InfoContextf(ctx context.Context, template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		InfoContext(ctx, fmt.Sprintf(template, args...))
		return
	}

	traceID := extractTraceID(ctx)
	// 使用调用者的位置创建临时logger，调整调用栈跳过层数
	infoLogger := logger.WithOptions(zap.AddCallerSkip(0))
	if traceID != DefaultTraceID {
		// 输出日志，添加trace_id
		infoLogger.Info(fmt.Sprintf(template, args...), zap.String("trace_id", traceID))
	} else {
		// 输出日志，不添加trace_id
		infoLogger.Info(fmt.Sprintf(template, args...))
	}
}

// WarnContextf 格式化输出带有trace_id的警告日志
func WarnContextf(ctx context.Context, template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		WarnContext(ctx, fmt.Sprintf(template, args...))
		return
	}

	traceID := extractTraceID(ctx)
	// 使用调用者的位置创建临时logger
	warnLogger := logger.WithOptions(zap.AddCallerSkip(0))
	if traceID != DefaultTraceID {
		// 输出日志，添加trace_id
		warnLogger.Warn(fmt.Sprintf(template, args...), zap.String("trace_id", traceID))
	} else {
		// 输出日志，不添加trace_id
		warnLogger.Warn(fmt.Sprintf(template, args...))
	}
}

// ErrorContextf 格式化输出带有trace_id的错误日志
func ErrorContextf(ctx context.Context, template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		ErrorContext(ctx, fmt.Sprintf(template, args...))
		return
	}

	traceID := extractTraceID(ctx)
	// 使用调用者的位置创建临时logger
	errorLogger := logger.WithOptions(zap.AddCallerSkip(0))
	if traceID != DefaultTraceID {
		// 输出日志，添加trace_id
		errorLogger.Error(fmt.Sprintf(template, args...), zap.String("trace_id", traceID))
	} else {
		// 输出日志，不添加trace_id
		errorLogger.Error(fmt.Sprintf(template, args...))
	}
}

// FatalContextf 格式化输出带有trace_id的致命错误日志
func FatalContextf(ctx context.Context, template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		FatalContext(ctx, fmt.Sprintf(template, args...))
		return
	}

	traceID := extractTraceID(ctx)
	// 使用调用者的位置创建临时logger
	fatalLogger := logger.WithOptions(zap.AddCallerSkip(0))
	if traceID != DefaultTraceID {
		// 输出日志，添加trace_id
		fatalLogger.Fatal(fmt.Sprintf(template, args...), zap.String("trace_id", traceID))
	} else {
		// 输出日志，不添加trace_id
		fatalLogger.Fatal(fmt.Sprintf(template, args...))
	}
}

// 保留原有方法不变

// Debug 输出调试日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 输出信息日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 输出警告日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 输出错误日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal 输出致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Debugf 格式化输出调试日志
func Debugf(template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		sugar.Debugf(template, args...)
		return
	}

	// 使用调用者的位置创建临时logger
	debugLogger := logger.WithOptions(zap.AddCallerSkip(0))

	// 输出日志，不添加额外字段，因为调用位置已经正确了
	debugLogger.Debug(fmt.Sprintf(template, args...))
}

// Infof 格式化输出信息日志
func Infof(template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		sugar.Infof(template, args...)
		return
	}

	// 使用调用者的位置创建临时logger，调整调用栈跳过层数
	infoLogger := logger.WithOptions(zap.AddCallerSkip(0))

	// 输出日志，位置信息会自动包含在日志中
	infoLogger.Info(fmt.Sprintf(template, args...))
}

// Warnf 格式化输出警告日志
func Warnf(template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		sugar.Warnf(template, args...)
		return
	}

	// 使用调用者的位置创建临时logger
	warnLogger := logger.WithOptions(zap.AddCallerSkip(0))

	// 输出日志，位置信息会自动包含在日志中
	warnLogger.Warn(fmt.Sprintf(template, args...))
}

// Errorf 格式化输出错误日志
func Errorf(template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		sugar.Errorf(template, args...)
		return
	}

	// 使用调用者的位置创建临时logger
	errorLogger := logger.WithOptions(zap.AddCallerSkip(0))

	// 输出日志，位置信息会自动包含在日志中
	errorLogger.Error(fmt.Sprintf(template, args...))
}

// Fatalf 格式化输出致命错误日志
func Fatalf(template string, args ...interface{}) {
	// 获取调用者信息
	_, _, _, ok := runtime.Caller(1)
	if !ok {
		// 如果获取调用者信息失败，使用标准糖化日志
		sugar.Fatalf(template, args...)
		return
	}

	// 使用调用者的位置创建临时logger
	fatalLogger := logger.WithOptions(zap.AddCallerSkip(0))

	// 输出日志，位置信息会自动包含在日志中
	fatalLogger.Fatal(fmt.Sprintf(template, args...))
}

// Sync 同步日志
func Sync() error {
	return logger.Sync()
}
