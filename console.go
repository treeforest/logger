package logger

import (
	"fmt"
	"time"
	"os"
)

// 往终端打印日志

// consoleLogger 终端日志打印结构体
type consoleLogger struct {
	level Level
}

func NewConsoleLogger(level Level) Logger {
	return &consoleLogger{
		level: level,
	}
}

// Debug 方法
func (c *consoleLogger) Debug(args ...interface{}) {
	c.log(DebugLevel, "%v", args...)
}
func (c *consoleLogger) Debugf(format string, args ...interface{}) {
	c.log(DebugLevel, format, args...)
}

// Info 方法
func (c *consoleLogger) Info(args ...interface{}) {
	c.log(InfoLevel, "%v", args...)
}
func (c *consoleLogger) Infof(format string, args ...interface{}) {
	c.log(InfoLevel, format, args...)
}

// Warn 方法
func (c *consoleLogger) Warn(args ...interface{}) {
	c.log(WarnLevel, "%v", args...)
}
func (c *consoleLogger) Warnf(format string, args ...interface{}) {
	c.log(WarnLevel, format, args...)
}

// Error 方法
func (c *consoleLogger) Error(args ...interface{}) {
	c.log(ErrorLevel, "%v", args...)
}
func (c *consoleLogger) Errorf(format string, args ...interface{}) {
	c.log(ErrorLevel, format, args...)
}

// Fatal 方法
func (c *consoleLogger) Fatal(args ...interface{}) {
	c.log(FatalLevel, "%v", args...)
}
func (c *consoleLogger) Fatalf(format string, args ...interface{}) {
	c.log(FatalLevel, format, args...)
}

func (c *consoleLogger) OpenDebug() {}
func (c *consoleLogger) CloseDebug() {}

// 设置日志级别
func (c *consoleLogger) SetLevel(level Level) {
	c.level = level
}

func (c *consoleLogger) Close() {
	//终端下的标准输出不用关闭
}

func (c *consoleLogger) log(level Level, format string, args ...interface{}) {
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
