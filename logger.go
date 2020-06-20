// 日志库文件
package log

var defaultLogger Logger

func SetFileLogger() {
	defaultLogger = NewFileLogger(1024*30, DebugLevel, "", "./")
	defaultLogger.(*fileLogger).skip = 4
}

func init() {
	defaultLogger = NewConsoleLogger(DebugLevel)
	defaultLogger.(*consoleLogger).skip = 4
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// Level 是一个uint16的自定义类型，代表日志级别
type Level uint16

// 定义的具体的日志级别常量
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

//日志级别对应的显示字符串
var levels = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

func getLevelStr(level Level) string {
	switch level {
	case DebugLevel:
		return levels[0]
	case InfoLevel:
		return levels[1]
	case WarnLevel:
		return levels[2]
	case ErrorLevel:
		return levels[3]
	case FatalLevel:
		return levels[4]
	default:
		return "DEBUG"
	}
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	OpenDebug()
	CloseDebug()

	SetLevel(level Level)
	Close()
}
