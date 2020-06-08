package main

import "github.com/treeforest/logger"

var GlobalLogger logger.Logger

func main() {
	GlobalLogger = logger.NewFileLogger(1 * 1024 * 1024, logger.DebugLevel, "./", "test.log")
	//GlobalLogger = logger.NewConsoleLogger(logger.DebugLevel)
	defer GlobalLogger.Close()

	for i :=0; i < 1024 * 20; i++ {
		GlobalLogger.Debug("%s 这是一条测试的日志。","Debug")
		GlobalLogger.Info("%s 这是一条测试的日志。","Info")
		GlobalLogger.Error("这是一条测试的日志。")
	}
}
