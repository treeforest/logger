// 日志库文件
package log

import (
	"fmt"
	"github.com/treeforest/logger/color"
	"time"
)

// 管道缓存日志单元的最大容量
const max_chan_size uint32 = 1000

// 检查文件过大的频率
const default_flush_tick = 10 * time.Second

// 默认全局日志记录句柄
var defaultLogger logger = newLogger(4)

//OnInit 初始化log配置
//@param : path: 路径，level: 过滤等级，
//		   size: 切割文件大小 jsonFile: 是否打印成json格式
func OnInit(path string, level LogLevel, size int64, jsonFile bool) {
	defaultLogger.(*loggerHandle).onInit(path, level, size, jsonFile)
}

func Debug(a ...interface{}) {
	defaultLogger.debug(a...)
}

func Debugf(format string, a ...interface{}) {
	defaultLogger.debug(fmt.Sprintf(format, a...))
}

func Info(a ...interface{}) {
	defaultLogger.info(a...)
}

func Infof(format string, a ...interface{}) {
	defaultLogger.info(fmt.Sprintf(format, a...))
}

func Warn(a ...interface{}) {
	defaultLogger.warn(a...)
}

func Warnf(format string, a ...interface{}) {
	defaultLogger.warn(fmt.Sprintf(format, a...))
}

func Error(a ...interface{}) {
	defaultLogger.error(a...)
}

func Errorf(format string, a ...interface{}) {
	defaultLogger.error(fmt.Sprintf(format, a))
}

func Fatal(a ...interface{}) {
	defaultLogger.fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	defaultLogger.fatal(fmt.Sprintf(format, a...))
}

func SetLogLevel(level LogLevel) {
	defaultLogger.setLogLevel(level)
}

// LogLevel 是一个uint16的自定义类型，代表日志级别
type LogLevel uint32

// 日志级别
const (
	LOGDEBUG LogLevel = iota
	LOGINFO
	LOGWARN
	LOGERROR
	LOGFATAL
)

// 日志打印颜色
var levelColors []*color.Color = []*color.Color{
	color.New(color.FgBlue),
	color.New(color.FgGreen),
	color.New(color.FgHiYellow),
	color.New(color.FgRed),
	color.New(color.FgRed),
}

func stdLogger(level LogLevel) *color.Color {
	switch level {
	case LOGDEBUG, LOGINFO, LOGWARN, LOGERROR, LOGFATAL:
		return levelColors[level]
	default:
		return levelColors[LOGDEBUG]
	}
}

// 日志级别对应的显示字符串
var levels = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

func getLevelStr(level LogLevel) string {
	switch level {
	case LOGDEBUG:
		return levels[LOGDEBUG]
	case LOGINFO:
		return levels[LOGINFO]
	case LOGWARN:
		return levels[LOGWARN]
	case LOGERROR:
		return levels[LOGERROR]
	case LOGFATAL:
		return levels[LOGFATAL]
	default:
		return levels[LOGDEBUG]
	}
}

type logger interface {
	debug(a ...interface{})
	info(a ...interface{})
	warn(a ...interface{})
	error(a ...interface{})
	fatal(a ...interface{})

	setLogLevel(level LogLevel)
	close()
}
