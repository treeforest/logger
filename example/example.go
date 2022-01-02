package main

import (
	"fmt"
	"time"

	log "github.com/treeforest/logger"
)

func main() {
	txt := "Hello, this is a test log entry"

	// log.SetLogger(log.NewFileLogger(".", 1024*1024*5))

	t := time.NewTimer(time.Second)
	count := 0
	for {
		select {
		case <-t.C:
			fmt.Println("count: ", count)
			return
		default:
			log.Info(txt)
			count++
		}
	}
}
