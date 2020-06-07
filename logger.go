
// 日志库文件
package logger

// Level 是一个uint16的自定义类型，代表日志级别
type Level uint16

// 具体的日志级别常量
const (
	DebugLevel Level= iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)