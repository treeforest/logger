
// 日志库文件
package logger

// Level 是一个uint16的自定义类型，代表日志级别
type Level uint16

// 定义的具体的日志级别常量
const (
	DebugLevel Level= iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

const DefaultFileSize int64 = 10 * 1024 * 1024

func getLevelStr(level Level) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "DEBUG"
	}
}