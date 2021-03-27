package main

import (
	"github.com/treeforest/logger"
	"sync"
)

func TestInGoroutine() {
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			log.Infof("Hello -> %d", j)
		}(i)
	}
	wg.Wait()
	log.Stop()
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

func TestPanic() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			log.Stop()
		}
	}()
	log.Debug("hello")
	panic("error")
}

func main() {
	//TestInGoroutine()
	//TestGetLogger()
	//TestDefaultLog()
	//TestWriteSuccess()
	//TestSetLevel()
	TestPanic()
}
