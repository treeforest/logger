package logger

import (
	"fmt"
	"time"
	"os"
)

// 往终端打印日志

// ConsoleLogger 终端日志打印结构体
type ConsoleLogger struct {
	level Level
}

func NewConsoleLogger(level Level) *ConsoleLogger {
	return &ConsoleLogger{
		level: level,
	}
}

// Debug 方法
func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	c.log(DebugLevel, format, args...)
}

// Info 方法
func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	c.log(InfoLevel, format, args...)
}

// Warn 方法
func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	c.log(WarningLevel, format, args...)
}

// Error 方法
func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	c.log(ErrorLevel, format, args...)
}

// Fatal 方法
func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	c.log(FatalLevel, format, args...)
}

// 设置日志级别
func (c *ConsoleLogger) SetLevel(level Level) {
	c.level = level
}

func (c *ConsoleLogger) log(level Level, format string, args ...interface{}) {
	if c.level > level {
		return
	}

	// 日志格式:[时间][文件:行号][函数名][日志级别] 日志信息
	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, funcName, line := getCallerInfo(3)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] %s", now, fileName, line, funcName, getLevelStr(level), msg)
	fmt.Fprintln(os.Stdout, logMsg)
}
