/*
 * 日志文件主要实现逻辑（写文件与控制台打印）
 */

package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 文件日志结构体信息
type loggerHandle struct {
	level               LogLevel       // 日志级别门槛，低于该级别的日志将不打印
	path                string         // 日志文件路径
	pathName            string         // 带有完整路径的文件名
	pFile               *os.File       // 存放一般的日志文件路径
	maxFileSize         int64          // 日志文件的最大大小
	initOnce            sync.Once      // 防止日志多次初始化
	jsonFile            bool           // 输出到文件的日志文件格式是否为json格式
	fileMutex           sync.Mutex     // 确保多协程读写文件，防止文件内容混乱，做到协程安全
	depth               uint32         // 调用函数的层级，默认为4
	localIP             string         // 本地ip
	closeWait           sync.WaitGroup // 等待结束
	onlyStdWriterCancel func()         // 只有控制台写时的上下文cancel
	logItemPool         sync.Pool      // 临时对象池
	stopChan            chan struct{}  // 停止打印信号
	logChan             chan *logItem  // 日志打印缓冲池
}

// 日志记录单元
type logItem struct {
	content   string
	level     LogLevel
	fileName  string // 文件名
	line      int    // 文件行号
	localFunc string // 本地函数名
	close     bool   // 如果调用fatal级别，该值将设为true，便于协程的退出
}

// 初始化配置
func newLogger(depth uint32) logger {
	var newLogItem = func() interface{} {
		return new(logItem)
	}

	l := &loggerHandle{
		depth:       depth,
		logItemPool: sync.Pool{New: newLogItem},
		logChan:     make(chan *logItem, max_chan_size),
		stopChan:    make(chan struct{}),
	}

	l.startStdLogger()

	return l
}

// Debug 方法
func (h *loggerHandle) debug(a ...interface{}) {
	h.log(LOGDEBUG, fmt.Sprint(a...))
}

func (h *loggerHandle) debugf(format string, a ...interface{}) {
	h.log(LOGDEBUG, fmt.Sprintf(format, a...))
}

// Info 方法
func (h *loggerHandle) info(a ...interface{}) {
	h.log(LOGINFO, fmt.Sprint(a...))
}

func (h *loggerHandle) infof(format string, a ...interface{}) {
	h.log(LOGINFO, fmt.Sprintf(format, a...))
}

// Warn 方法
func (h *loggerHandle) warn(a ...interface{}) {
	h.log(LOGWARN, fmt.Sprint(a...))
}

func (h *loggerHandle) warnf(format string, a ...interface{}) {
	h.log(LOGWARN, fmt.Sprintf(format, a...))
}

// Error 方法
func (h *loggerHandle) error(a ...interface{}) {
	h.log(LOGERROR, fmt.Sprint(a...))
}

func (h *loggerHandle) errorf(format string, a ...interface{}) {
	h.log(LOGERROR, fmt.Sprintf(format, a...))
}

// Fatal 方法
func (h *loggerHandle) fatal(a ...interface{}) {
	h.log(LOGFATAL, fmt.Sprint(a...))
	log.Println("!!! log call fatal exit !!!\n")
	os.Exit(1)
}

func (h *loggerHandle) fatalf(format string, a ...interface{}) {
	h.log(LOGFATAL, fmt.Sprintf(format, a...))
	log.Println("!!! log call fatalf exit !!!\n")
	os.Exit(1)
}

// 设置日志级别
func (h *loggerHandle) setLogLevel(level LogLevel) {
	h.level = level
}

// 清理资源
func (h *loggerHandle) close() {
	h.closeWait.Wait()

	if h.stopChan != nil {
		close(h.stopChan)
	}

	if h.logChan != nil {
		close(h.logChan)
	}

	if h.pFile != nil {
		h.pFile.Close()
		h.pFile = nil
	}
}

// 启动控制台日志文件打印
func (h *loggerHandle) startStdLogger() {
	if h.logChan == nil {
		h.logChan = make(chan *logItem, max_chan_size)
	}

	if h.stopChan == nil {
		h.stopChan = make(chan struct{})
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.onlyStdWriterCancel = cancel
	go h.onlyStdWriter(ctx)
}

// 仅开启控制台进行日志打印
func (h *loggerHandle) onlyStdWriter(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("onlyStdWriter() panic: %v\n", err)
		}
		h.closeWait.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-h.logChan:
			if ok {
				toStdStr, _ := h.unpack(item)
				h.outputConsole(item.level, &toStdStr)
				if item.close {
					close(h.stopChan)
					return
				}
			} else {
				// 将缓存的日志执行结束
				for item := range h.logChan {
					toStdStr, _ := h.unpack(item)
					h.outputConsole(item.level, &toStdStr)
				}

				panic("logChan has closed!!!")
			}
			h.putLogItem(item)
		}
	}
}

// 同时开启控制台与文件的日志记录
func (h *loggerHandle) logWriter() {
	defer func() {
		if h.pFile != nil {
			h.pFile.Close()
		}
		if err := recover(); err != nil {
			log.Printf("logWriter() panic: %v\n", err)
		}
		h.closeWait.Done()
	}()

	// 检测文件大小的定时器，如果文件超过设定的阈值，则进行切分文件
	flushTicker := time.NewTicker(default_flush_tick)
	defer flushTicker.Stop()

	// 定时切分文件的定时器，区别每天的文件
	switchTimer := time.NewTimer(getFirstSwitchTime())
	defer switchTimer.Stop()

	for {
		select {
		case item, ok := <-h.logChan:
			if ok {
				toStdStr, toFileStr := h.unpack(item)
				h.outputFile(&toFileStr)
				h.outputConsole(item.level, &toStdStr)
				if item.close {
					close(h.stopChan)
					return
				}
			} else {
				for item := range h.logChan {
					toStdStr, toFileStr := h.unpack(item)
					h.outputFile(&toFileStr)
					h.outputConsole(item.level, &toStdStr)
				}

				panic("logChan has closed!!!")
			}
			h.putLogItem(item)

		case <-flushTicker.C:
			if fileSize(h.pathName) >= h.maxFileSize {
				h.fileMutex.Lock()
				if err := h.switchFile(); err != nil {
					log.Printf("log switch failed: %v\n", err)
					h.fileMutex.Unlock()
					panic(err)
				}
				h.fileMutex.Unlock()
			}

		case <-switchTimer.C:
			h.fileMutex.Lock()
			if err := h.switchFile(); err != nil {
				log.Printf("log switch failed: %v\n", err)
				h.fileMutex.Unlock()
				panic(err)
			}
			h.fileMutex.Unlock()
			switchTimer.Reset(getNextSwitchTime())
		}
	}
}

// 输出到控制台
func (h *loggerHandle) outputConsole(level LogLevel, s *string) {
	stdLogger(level).Println(*s)
}

// 输出到日志文件
func (h *loggerHandle) outputFile(s *string) {
	if h.pFile != nil {
		fmt.Fprintln(h.pFile, *s)
	}
}

// 初始化配置
func (h *loggerHandle) onInit(path string, level LogLevel, size int64, jsonFile bool) {
	setupFunc := func() {
		h.path = rectifyPath(path)
		h.level = level
		h.depth = 4
		h.pFile = nil
		h.maxFileSize = size
		h.jsonFile = jsonFile
		h.localIP = localIP()
		if h.logChan == nil {
			h.logChan = make(chan *logItem, max_chan_size)
		}
		if h.stopChan == nil {
			h.stopChan = make(chan struct{})
		}

		if err := h.switchFile(); err != nil {
			log.Printf("switchFile() panic: %v\n", err)
			panic(err)
		}

		h.onlyStdWriterCancel()

		h.closeWait.Add(1)
		go h.logWriter()
	}

	h.initOnce.Do(setupFunc)
}

// 记录
func (h *loggerHandle) log(level LogLevel, content string) {
	if h.level > level {
		return
	}

	// 日志格式:[日志级别][时间][文件:行号]: 日志信息
	item := h.getLogItem()
	item.level = level
	item.content = content
	item.fileName, item.localFunc, item.line = getCallerInfo(h.depth)
	if h.level == LOGFATAL {
		item.close = true
	} else {
		item.close = false
	}

	h.logChan <- item
}

// 解包操作
func (h *loggerHandle) unpack(item *logItem) (toStd string, toFile string) {
	now := time.Now().Format("2006-01-02 15:04:05.000")

	toStd = fmt.Sprintf("[%s][%s:%d][%s]: %s", now, item.fileName, item.line, getLevelStr(item.level), item.content)
	if h.jsonFile {
		toFile = fmt.Sprintf("{\"LEVEL\":\"%s\",\"Time\":\"%v\",\"File\":\"%s\",\"Line\":\"%s\",\"LocalFunc\":\"%s\",\"CONTENT\":%s}",
			item.level, now, item.fileName, item.line, item.localFunc, item.content)
	} else {
		toFile = toStd
	}
	return
}

// 放回临时对象池
func (h *loggerHandle) putLogItem(l *logItem) {
	h.logItemPool.Put(l)
}

// 获取log对象
func (h *loggerHandle) getLogItem() *logItem {
	return h.logItemPool.Get().(*logItem)
}

//取基础文件名
func (h *loggerHandle) getFileName() (filename string) {
	now := time.Now()
	filename = h.path + now.Format("20060102") + "." + h.localIP //+ ".log"
	return
}

// 切换文件
func (h *loggerHandle) switchFile() error {
	fileName := h.getFileName()

	// 确认目录存在
	if err := os.MkdirAll(h.path, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	// 先关闭旧文件再切换
	if h.pFile != nil {
		if err := h.pFile.Close(); err != nil {
			return err
		}
		h.pFile = nil
	}

	// 创建或者打开已存在文件
	file, err := h.newFile(fileName)
	if err != nil {
		return err
	}
	h.pFile = file

	return nil
}

// 新建文件，返回描叙符
func (h *loggerHandle) newFile(fileName string) (*os.File, error) {

	h.pathName = fileName + ".log"

	for fileID := 2; fileExists(h.pathName); fileID++ {
		h.pathName = fileName + fmt.Sprintf(".%02d.log", fileID)
	}

	f, err := os.OpenFile(h.pathName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	//log.Printf("newFile pathname %s \n", pathname)
	return f, nil
}
