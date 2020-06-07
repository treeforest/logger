package logger

import (
	"os"
	"path"
	"fmt"
	"time"
	)

// 往文件里面写日志

// 文件日志结构体信息
type FileLogger struct {
	level Level			// 日志级别门槛，低于该级别的日志将不打印
	fileName string
	filePath string
	file *os.File
	errFile *os.File
	maxSize int64
}

func NewFileLogger(level Level, fileName, filePath string) *FileLogger {
	fileLogger := &FileLogger{
		level: level,
		fileName: fileName,
		filePath: filePath,
		maxSize: DefaultFileSize,
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

func (f *FileLogger) log(level Level, format string, args ...interface{}) {
	if f.level > level {
		return
	}

	// 日志格式:[时间][文件:行号][函数名][日志级别] 日志信息
	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, funcName, line := getCallerInfo(3)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] %s", now, fileName, line, funcName, getLevelStr(level), msg)

	if f.checkSplit(f.file) {
		f.file = f.splitLogFile(f.file)
	}
	fmt.Fprintln(f.file, logMsg)

	// 如果是Error或者Fatal级别的日志还要记录到 f.errFile
	if level >= ErrorLevel {
		if f.checkSplit(f.errFile) {
			f.errFile = f.splitLogFile(f.errFile)
		}
		fmt.Fprintln(f.errFile, logMsg)
	}
}

func (f *FileLogger) checkSplit(file *os.File) bool {
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	return fileSize >= f.maxSize
}

func (f *FileLogger) splitLogFile(file *os.File) *os.File{
	fileName := file.Name()
	// 切分文件
	backupName := fmt.Sprintf("%s_%v.back", fileName, time.Now().Unix())
	// 1. 把原来的文件关闭
	file.Close()
	// 2. 备份原来的文件
	os.Rename(fileName, backupName)
	// 3. 新建一个文件
	newFile, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open logfile(%s) failed. err:%v", fileName, err))
	}
	return newFile
}

// Debug 方法
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.log(DebugLevel, format, args...)
}

// Info 方法
func (f *FileLogger) Info(format string, args ...interface{}) {
	f.log(InfoLevel, format, args...)
}

// Warn 方法
func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.log(WarningLevel, format, args...)
}

// Error 方法
func (f *FileLogger) Error(format string, args ...interface{}) {
	f.log(ErrorLevel, format, args...)
}

// Fatal 方法
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.log(FatalLevel, format, args...)
}

// 设置日志级别
func (f *FileLogger) SetLevel(level Level) {
	f.level = level
}

func (f *FileLogger) SetMaxSize(maxSize int64) {
	f.maxSize = maxSize
}