package main

import "github.com/treeforest/logger"

func main() {
	//l := logger.NewFileLogger(logger.DebugLevel, "./", "test.log")
	l := logger.NewConsoleLogger(logger.DebugLevel)

	l.Debug("%s 这是一条测试的日志。","Debug")
	l.Info("%s 这是一条测试的日志。","Info")
	l.Error("这是一条测试的日志。")
}
