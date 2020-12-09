package main

import (
	"github.com/treeforest/logger"
)

func TestGetLogger() {
	logger := log.GetLogger("log", log.WithFilePath("./log/"),
		log.WithErrFilePath("./log_error/"),
		log.WithJsonFile(true),
		log.WithLogLevel(log.DebugLevel))

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

func main() {
	//TestGetLogger()
	TestDefaultLog()
	//TestWriteSuccess()
}
