package log

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

// 存放公共的工具函数
func getCallerInfo(skip int) (fileName, funcName string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	// 从file(x/y/xx.go)中获取文件名
	fileName = path.Base(file)
	// 根据pc拿到函数名
	funcName = path.Base(runtime.FuncForPC(pc).Name())
	return
}

// 获取日志文件名，以时间为记录节点
func getFileLoggerNameByTime() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	min := now.Minute()
	sec := now.Second()
	var filename string
	filename = fmt.Sprintf("%d", year)
	if month < 10 {
		filename = filename + fmt.Sprintf("0%d", month)
	} else {
		filename = filename + fmt.Sprintf("%d", month)
	}
	if day < 10 {
		filename = filename + fmt.Sprintf("0%d", day)
	} else {
		filename = filename + fmt.Sprintf("%d", day)
	}
	if hour < 10 {
		filename = filename + fmt.Sprintf("0%d", hour)
	} else {
		filename = filename + fmt.Sprintf("%d", hour)
	}
	if min < 10 {
		filename = filename + fmt.Sprintf("0%d", min)
	} else {
		filename = filename + fmt.Sprintf("%d", min)
	}
	if sec < 10 {
		filename = filename + fmt.Sprintf("0%d", sec)
	} else {
		filename = filename + fmt.Sprintf("%d", sec)
	}

	return filename + ".log"
}
