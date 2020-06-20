package main

import (
	"time"
	"github.com/treeforest/logger"
)

func main() {

	go func() {
		for i := 0; i < 20; i++ {
			log.Debug("----这是一条测试的日志----")
			log.Info("----这是一条测试的日志----")
			log.Warn("----这是一条测试的日志----")
			log.Error("----这是一条测试的日志----")
			log.Fatal("----这是一条测试的日志----")
		}
	}()

	// log.SetFileLogger()

	//for i := 0; i < 20; i++ {
	//	log.Debug("----这是一条测试的日志----")
	//	log.Info("----这是一条测试的日志----")
	//	log.Warn("----这是一条测试的日志----")
	//	log.Error("----这是一条测试的日志----")
	//	log.Fatal("----这是一条测试的日志----")
	//}

	time.Sleep(1 * time.Second)
	//for i :=0; i < 1024 * 20; i++ {
	//	GlobalLogger.Debug("%s 这是一条测试的日志。","Debug")
	//	GlobalLogger.Info("%s 这是一条测试的日志。","Info")
	//	GlobalLogger.Error("这是一条测试的日志。")
	//}
}
