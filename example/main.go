package main

import (
	"fmt"
	"github.com/treeforest/logger"
	"time"
)

func TestInGoroutine() {
	now := time.Now()
	for i := 0; i < 100000; i++ {
		log.Infof("Hello -> %d", i)
		//fmt.Println("----------")
	}
	interval := time.Now().Sub(now).Milliseconds()
	fmt.Printf("时间: %dms", interval)
}

func TestGetLogger() {
	logger := log.GetLogger("log", log.WithFilePath("./log/"),
		log.WithErrFilePath("./log_error/"),
		log.WithJsonFile(true),
		log.WithLogLevel(log.DebugLevel))
	defer logger.Stop()

	logger.Debug("Debug Message")
	logger.Info("Info Message")
	logger.Warn("Warn Message")
	logger.Error("Error Message")
	logger.Fatal("Fatal Message...")
}

func TestDefaultLog() {
	defer log.Stop()

	log.SetConfig(
		log.WithLogLevel(log.InfoLevel),
		log.WithFilePath("."))

	for i := 0; i < 100; i++ {
		log.Debug("Debug Message")
		log.Info("Info Message")
		log.Warn("Warn Message")
		log.Error("Error Message")
	}
}

func TestWriteSuccess() {
	defer log.Stop()

	for i := 0; i < 10; i++ {
		log.Debug("Debug Message")
		log.Info("Info Message")
		log.Warn("Warn Message")
		log.Error("Error Message")
	}
}

func TestSetLevel() {
	defer log.Stop()
	log.SetConfig(log.WithLogLevel(log.InfoLevel))
	log.Debug("Hello Debug")
	log.Info("Hello Info")
}

func main() {
	TestInGoroutine()
	//TestGetLogger()
	//TestDefaultLog()
	//TestWriteSuccess()
	//TestSetLevel()
}
