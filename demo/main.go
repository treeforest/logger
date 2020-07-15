package main

import (
	"time"
	"github.com/treeforest/logger"
	"fmt"
)

func main() {

	log.OnInit("./", log.LOGDEBUG, 1024 * 3, false)
	now := time.Now()
	for i := 0; i < 100; i++ {
		log.Debug("----这是一条 Debug 的日志----")
		log.Info("----这是一条 Info 的日志----")
		log.Warn("----这是一条 Warn 的日志----")
		log.Error("----这是一条 Error 的日志----")
	}

	log.Fatal("----这是一条 Fatal 的日志----")

	sub := time.Now().Sub(now).Seconds()

	fmt.Println(sub)

}
