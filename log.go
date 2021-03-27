// 日志库文件
package log

import (
	"sync"
)

func init() {
	defaultLogger = GetLogger(defaultNoneModule, WithLogLevel(DebugLevel))
}

// 默认全局日志记录句柄
var defaultLogger Logger

func Debug(a ...interface{}) {
	defaultLogger.Debug(a...)
}

func Debugf(format string, a ...interface{}) {
	defaultLogger.Debugf(format, a...)
}

func Info(a ...interface{}) {
	defaultLogger.Info(a...)
}

func Infof(format string, a ...interface{}) {
	defaultLogger.Infof(format, a...)
}

func Warn(a ...interface{}) {
	defaultLogger.Warn(a...)
}

func Warnf(format string, a ...interface{}) {
	defaultLogger.Warnf(format, a...)
}

func Error(a ...interface{}) {
	defaultLogger.Error(a...)
}

func Errorf(format string, a ...interface{}) {
	defaultLogger.Errorf(format, a...)
}

func Fatal(a ...interface{}) {
	defaultLogger.Fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	defaultLogger.Fatalf(format, a...)
}

func SetConfig(opts ...Option) {
	defaultLogger.SetConfig(opts...)
}

// Stop 关闭日志。用户主动关闭日志读写，否则会在异步情况下丢失日志数据
func Stop() {
	defaultLogger.Stop()
}

// logLevel 是一个uint16的自定义类型，代表日志级别
type logLevel uint32

// 日志级别
const (
	DebugLevel logLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})

	SetConfig(opts ...Option)
	Stop()
}

var loggers = make(map[string]Logger)
var lock = sync.Mutex{}

// GetLogger 获取不同模块下的日志句柄，若改模块已存在，则直接返回句柄；
// 否则，创建一个新的日志句柄并返回。
func GetLogger(module string, opts ...Option) Logger {
	lock.Lock()
	defer lock.Unlock()

	if module == "" {
		module = defaultNoneModule
	}

	if l, ok := loggers[module]; ok {
		return l
	}

	l := newLogger(append(opts, WithModule(module))...)
	loggers[module] = l
	return l
}

// StopAll 停止所有的日志模块。在检测到panic时，建议停止所有
// 的日志模块，以确保异步的日志能在panic之前打印结束。
func StopAll() {
	for _, l := range loggers {
		l.Stop()
	}
}