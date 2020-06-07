package logger

import (
	"os"
	"path"
	"fmt"
	"time"
)

// 往文件里面写日志

// FileLogger struct
type FileLogger struct {
	fileName string
	filePath string
	file *os.File
	errFile *os.File
}

func NewFileLogger(fileName, filePath string) *FileLogger {
	fileLogger := &FileLogger{
		fileName: fileName,
		filePath: filePath,
	}
	fileLogger.initFile()
	return fileLogger
}

func (f *FileLogger) initFile() {
	logName := path.Join(f.filePath, f.fileName)

	// open file
	file, err := os.OpenFile(logName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open logfile(%s) failed. err:%v", logName, err))
	}
	f.file = file

	// open error file
	errLogName := fmt.Sprintf("%s.error", logName)
	errFile, err := os.OpenFile(errLogName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open logfile(%s) failed. err:%v", errLogName, err))
	}
	f.errFile = errFile
}

//
func (f *FileLogger) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args)
	// 日志格式:[时间][文件:行号][函数名][日志级别] 日志信息
	now := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, line, funcName := getCallerInfo(2)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] %s", now, fileName, line, funcName, "debug", msg)
	fmt.Fprintln(f.file, logMsg)
}