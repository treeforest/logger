package logger

import "github.com/fatih/color"

var defaultLogger Logger

func init() {
	// 初始化默认为控制台输出
	l := NewStdLogger(WithLogLevel(DEBUG))
	SetLogger(l)
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	SetLevel(lvl Level)

	Stop()
}

type Level int

const (
	DEBUG Level = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
)

var mapping map[Level]string = map[Level]string{
	DEBUG: "DEBU",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERRO",
	FATAL: "FATA",
}

var colorMapping map[Level]string = map[Level]string{
	DEBUG: "DEBU",
	INFO:  color.GreenString("INFO"),
	WARN:  color.YellowString("WARN"),
	ERROR: color.RedString("ERRO"),
	FATAL: color.MagentaString("FATA"),
}

// SetLogger 设置默认的日志对象
func SetLogger(logger Logger) {
	switch logger.(type) {
	case *stdLogger:
		logger.(*stdLogger).callDepth = 4
	case *fileLogger:
		logger.(*fileLogger).callDepth = 3
	}
	defaultLogger = logger
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func SetLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

func Stop() {
	defaultLogger.Stop()
}
