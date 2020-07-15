package main

import (
	"time"
	"github.com/treeforest/logger"
	"fmt"
)

func main() {

	log.OnInit("./", log.LOGDEBUG, 1024*1024*30, false)
	now := time.Now()
	for i := 0; i < 200; i++ {
		log.Debug("----这是一条测试的日志----")
		log.Info("----这是一条测试的日志----")
		log.Warn("----这是一条测试的日志----")
		log.Error("----这是一条测试的日志----")
		//log.Fatal("----这是一条测试的日志----")
	}

	sub := time.Now().Sub(now).Seconds()

	time.Sleep(time.Second * 2)

	fmt.Println(sub)

	//for i :=0; i < 1024 * 20; i++ {
	//	GlobalLogger.Debug("%s 这是一条测试的日志。","Debug")
	//	GlobalLogger.Info("%s 这是一条测试的日志。","Info")
	//	GlobalLogger.Error("这是一条测试的日志。")
	//}
}
