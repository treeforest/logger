package log

import (
	"fmt"
	"github.com/treeforest/logger/color"
	"os"
	"path"
	"sync"
	"time"
)

// 往文件里面写日志

// 文件日志结构体信息
type fileLogger struct {
	level       Level        // 日志级别门槛，低于该级别的日志将不打印
	fileName    string       // 日志文件名
	fileNameNil bool         // 初始化的时候，日志文件名如果为空，则默认使用当前创建的时间作为日志文件名
	filePath    string       // 日志文件路径
	file        *os.File     // 存放一般的日志文件路径
	errFile     *os.File     // 存放error日志的文件
	maxSize     int64        // 日志文件的最大大小
	mu          sync.Mutex   // 确保多协程读写文件，防止文件内容混乱，做到协程安全
	debugClose  bool         // 是否打印调试debug信息
	skip        int          // 调用函数的层级，默认为3
	debugColor  *color.Color // Debug 终端打印颜色
	infoColor   *color.Color // Info 终端打印颜色
	warnColor   *color.Color // Warn 终端打印颜色
	errorColor  *color.Color // Error 终端打印颜色
	fatalColor  *color.Color // Fatal 终端打印颜色
}

func NewFileLogger(maxSize int64, level Level, fileName, filePath string) Logger {
	fileLogger := &fileLogger{
		level:      level,
		fileName:   fileName,
		filePath:   filePath,
		maxSize:    maxSize,
		debugClose: false,
		skip:       3,
		debugColor: color.New(color.FgBlue),
		infoColor:  color.New(color.FgGreen),
		warnColor:  color.New(color.FgHiYellow),
		errorColor: color.New(color.FgRed),
		fatalColor: color.New(color.FgRed),
	}

	if fileName == "" {
		fileLogger.fileNameNil = true
		fileLogger.fileName = getFileLoggerNameByTime()
	}

	fileLogger.initFile()
	return fileLogger
}

// Debug 方法
func (f *fileLogger) Debug(args ...interface{}) {
	f.log(DebugLevel, "%v", args...)
}
func (f *fileLogger) Debugf(format string, args ...interface{}) {
	f.log(DebugLevel, format, args...)
}

// Info 方法
func (f *fileLogger) Info(args ...interface{}) {
	f.log(InfoLevel, "%v", args...)
}
func (f *fileLogger) Infof(format string, args ...interface{}) {
	f.log(InfoLevel, format, args...)
}

// Warn 方法
func (f *fileLogger) Warn(args ...interface{}) {
	f.log(WarnLevel, "%v", args...)
}
func (f *fileLogger) Warnf(format string, args ...interface{}) {
	f.log(WarnLevel, format, args...)
}

// Error 方法
func (f *fileLogger) Error(args ...interface{}) {
	f.log(ErrorLevel, "%v", args...)
}
func (f *fileLogger) Errorf(format string, args ...interface{}) {
	f.log(ErrorLevel, format, args...)
}

// Fatal 方法
func (f *fileLogger) Fatal(args ...interface{}) {
	f.log(FatalLevel, "%v", args...)
}
func (f *fileLogger) Fatalf(format string, args ...interface{}) {
	f.log(FatalLevel, format, args...)
}

// 是否开启debug时日志输出
func (f *fileLogger) OpenDebug() {
	f.debugClose = false
}

func (f *fileLogger) CloseDebug() {
	f.debugClose = true
}

// 设置日志级别
func (f *fileLogger) SetLevel(level Level) {
	f.level = level
}

// 关闭文件句柄
func (f *fileLogger) Close() {
	f.file.Close()
	f.errFile.Close()
}

func (f *fileLogger) initFile() {
	logName := path.Join(f.filePath, f.fileName)

	// open file
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open logfile(%s) failed. err:%v", logName, err))
	}
	f.file = file

	// open error file
	errLogName := fmt.Sprintf("%s.error", logName)
	errFile, err := os.OpenFile(errLogName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open logfile(%s) failed. err:%v", errLogName, err))
	}
	f.errFile = errFile
}

func (f *fileLogger) log(level Level, format string, args ...interface{}) {
	if f.level > level {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 日志格式:[时间][文件:行号][日志级别] 日志信息
	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, _, line := getCallerInfo(f.skip)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s] %s", now, fileName, line, getLevelStr(level), msg)

	if f.checkSplit(f.file) {
		f.file = f.splitLogFile(f.file)
	}

	// 写入文件
	fmt.Fprintln(f.file, logMsg)

	// 终端打印
	if !f.debugClose {
		switch level {
		case DebugLevel:
			f.debugColor.Println(logMsg)
		case InfoLevel:
			f.infoColor.Println(logMsg)
		case WarnLevel:
			f.warnColor.Println(logMsg)
		case ErrorLevel:
			f.errorColor.Println(logMsg)
		case FatalLevel:
			f.fatalColor.Println(logMsg)
		}
	}

	// 如果是Error或者Fatal级别的日志还要记录到 f.errFile
	if level >= ErrorLevel {
		if f.checkSplit(f.errFile) {
			f.errFile = f.splitLogFile(f.errFile)
		}
		fmt.Fprintln(f.errFile, logMsg)
	}
}

func (f *fileLogger) checkSplit(file *os.File) bool {
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	return fileSize >= f.maxSize
}

func (f *fileLogger) splitLogFile(file *os.File) *os.File {
	var newFile *os.File
	var err error

	if file == f.file {
		if f.fileNameNil {
			f.fileName = getFileLoggerNameByTime()
		} else {
			f.fileName = file.Name()
		}
	} else {
		if f.fileNameNil {
			f.fileName = getFileLoggerNameByTime() + ".error"
		} else {
			f.fileName = file.Name()
		}
	}

	if f.fileNameNil {
		file.Close()
		newFile, err = os.OpenFile(f.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			panic(fmt.Errorf("open logfile(%s) failed. err:%v", f.fileName, err))
		}
	} else {
		// 切分文件
		backupName := fmt.Sprintf("%s_%v.back", f.fileName, time.Now().Unix())
		// 1. 把原来的文件关闭
		file.Close()
		// 2. 备份原来的文件
		os.Rename(f.fileName, backupName)
		// 3. 新建一个文件
		newFile, err = os.OpenFile(f.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			panic(fmt.Errorf("open logfile(%s) failed. err:%v", f.fileName, err))
		}
	}

	return newFile
}
