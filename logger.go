
// 日志库文件
package logger

var CLog Logger = NewConsoleLogger(FatalLevel)

// Level 是一个uint16的自定义类型，代表日志级别
type Level uint16

// 定义的具体的日志级别常量
const (
	DebugLevel Level= iota
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