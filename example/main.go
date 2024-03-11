package main

import (
	"fmt"
	"time"

	log "github.com/treeforest/logger"
)

const (
	txt = "Hello, this is a test log entry"
)

func run(l log.Logger) int {
	defer l.Stop()

	t := time.NewTimer(time.Second)
	count := 0
	for {
		select {
		case <-t.C:
			fmt.Println("count: ", count)
			return count
		default:
			l.Debug(txt)
			count++
		}
	}
}

func main() {
	log.Info("welcome to use logger")

	fn := func(l log.Logger) {
		since := time.Now()
		count := run(l)
		used := time.Now().Sub(since)
		fmt.Printf("used:%dms average:%dns\n", used.Milliseconds(), used.Nanoseconds()/int64(count))
	}

	// 文件异步写
	fn(log.NewAsyncFileLogger(
		".",
		1024*1024*8,
		1024*64,
		time.Second,
		log.WithLogLevel(log.DEBUG),
		log.WithPrefix("example"),
	))

	// 文件同步写
	fn(log.NewSyncFileLogger(
		".",
		1024*1024*8,
		log.WithLogLevel(log.DEBUG),
		log.WithPrefix("example"),
	))

	// 控制台输出
	// fn(log.NewStdLogger())
}
