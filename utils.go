package logger

import (
	"runtime"
	"path"
	)

// 存放公共的工具函数

func getCallerInfo(skip int) (fileName, funcName  string, line int){
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