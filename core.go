/*
 * 日志文件主要实现逻辑（写文件与控制台打印）
 */

package log

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

//type core interface {
//	debug(a ...interface{})
//	info(a ...interface{})
//	warn(a ...interface{})
//	error(a ...interface{})
//	fatal(a ...interface{})
//}

var (
	RunState int32 = 0
	StoppingState int32 = 1
	StoppedState int32 = 2
)

// 文件日志结构体信息
type loggerCore struct {
	config          *Config       // 配置
	pathName        string        // 带有完整路径的文件名
	errPathName     string        // 带有完整路径的错误文件名
	pFile           *os.File      // 存放一般的日志文件句柄
	pErrFile        *os.File      // 存放错误的日志句柄
	fileWriter      *bufio.Writer // 日志文件缓冲写
	errFileWriter   *bufio.Writer // 错误日志文件缓冲写
	initOnce        sync.Once     // 防止日志多次初始化
	depth           uint32        // 调用函数的层级，默认为4
	localIP         string        // 本地ip
	cancelStdWriter func()        // 只有控制台写时的上下文cancel
	state           int32         // 运行状态
	stopped         chan bool     // 停止
	fatalStop       chan bool     // fatal日志导致的stop
	pool            sync.Pool     // 临时对象池
	items           chan *logItem // 日志打印缓冲池
}

// 日志记录单元
type logItem struct {
	content   string
	level     logLevel
	fileName  string // 文件名
	line      int    // 文件行号
	localFunc string // 本地函数名
}

// 初始化配置
func newLoggerCore(depth uint32, opts ...Option) *loggerCore {
	config := new(Config)
	for _, opt := range opts {
		opt(config)
	}
	config.validateConfig()

	var newLogItem = func() interface{} {
		return new(logItem)
	}

	core := &loggerCore{
		config:          config,
		depth:           depth,
		pool:            sync.Pool{New: newLogItem},
		items:           make(chan *logItem, config.maxChannelSize),
		state:           RunState,
		stopped:         make(chan bool),
		fatalStop:       make(chan bool),
		cancelStdWriter: func() {},
	}

	if core.config.path == "" {
		//log.Println("start std logger...")
		core.startStdLogger()
	} else {
		//log.Println("start std and file logger...")
		core.startStdAndFileLogger()
	}

	return core
}

func (lc *loggerCore) setConfig(opts ...Option) {
	for _, opt := range opts {
		opt(lc.config)
	}

	lc.config.validateConfig()

	if lc.config.path != "" {
		lc.startStdAndFileLogger()
	}
}

// 启动控制台日志文件打印
func (lc *loggerCore) startStdLogger() {
	ctx, cancel := context.WithCancel(context.Background())
	lc.cancelStdWriter = cancel
	go lc.stdWriter(ctx)
}

// 启动日志文件输出
func (lc *loggerCore) startStdAndFileLogger() {
	lc.initOnce.Do(func() {
		lc.pFile = nil
		lc.pErrFile = nil
		lc.localIP = localIP()

		if err := lc.switchFile(); err != nil {
			log.Printf("switchFile() panic: %v\n", err)
			panic(err)
		}

		if err := lc.switchErrFile(); err != nil {
			log.Printf("switchErrFile() panic: %v\n", err)
			panic(err)
		}

		// 关闭之前运行的控制台打印（若是开启）
		lc.cancelStdWriter()

		go lc.stdAndFileWriter()
	})
}

// debug 方法
func (lc *loggerCore) debug(a ...interface{}) {
	lc.log(DebugLevel, fmt.Sprint(a...))
}

// info 方法
func (lc *loggerCore) info(a ...interface{}) {
	lc.log(InfoLevel, fmt.Sprint(a...))
}

// warn 方法
func (lc *loggerCore) warn(a ...interface{}) {
	lc.log(WarnLevel, fmt.Sprint(a...))
}

// error 方法
func (lc *loggerCore) error(a ...interface{}) {
	lc.log(ErrorLevel, fmt.Sprint(a...))
}

// fatal 方法
func (lc *loggerCore) fatal(a ...interface{}) {
	lc.log(FatalLevel, fmt.Sprint(a...))

	<-lc.fatalStop
	//log.Println("!!! log call fatal exit !!!")
	os.Exit(1)
}

// 记录
func (lc *loggerCore) log(level logLevel, content string) {
	if lc.config.level > level || lc.state != RunState {
		return
	}

	item := lc.get()
	item.level = level
	item.content = content
	item.fileName, item.localFunc, item.line = getCallerInfo(lc.depth)

	// 异步写
	lc.items <- item
}

// 设置日志级别
func (lc *loggerCore) setLogLevel(level logLevel) {
	lc.config.level = level
}

// 仅开启控制台进行日志打印
func (lc *loggerCore) stdWriter(ctx context.Context) {
	isFatalLog := false

	defer func() {
		if err := recover(); err != nil {
			log.Printf("stdWriter() panic: %v\n", err)
		}
		if isFatalLog {
			lc.fatalStop <- true
		} else if atomic.LoadInt32(&lc.state) == StoppingState {
			//log.Println("stdWriter stopped.")
			lc.stopped <- true
		} else {
			// log.Println("cancel()")
		}
		//log.Println("exit...")
	}()
	//log.Println("start writer...")

	for {
		select {
		case <-ctx.Done():
			return

		case item, ok := <-lc.items:
			if !ok {
				//log.Println("channel has closed.")
				return
			}

			toStdStr, _ := lc.unpack(item)
			lc.outputConsole(item.level, &toStdStr)

			if item.level == FatalLevel {
				isFatalLog = true// 控制台输出是顺序性的，直接返回即可。若有缓存的日志，直接抛弃。
				return
			}
			lc.put(item)
		}
	}
}

// 同时开启控制台与文件的日志记录
func (lc *loggerCore) stdAndFileWriter() {
	// 检测文件大小的定时器，如果文件超过设定的阈值，则进行切分文件
	flushTicker := time.NewTicker(lc.config.fileFlushTick)

	// 定时切分文件的定时器，区别每天的文件
	switchTimer := time.NewTimer(getFirstSwitchTime())

	isFatalLog := false

	defer func() {
		lc.fileWriter.Flush()
		lc.errFileWriter.Flush()

		if err := recover(); err != nil {
			log.Printf("stdAndFileWriter() panic: %v\n", err)
		}

		flushTicker.Stop()
		switchTimer.Stop()

		if isFatalLog {
			lc.fatalStop <- true
		} else if atomic.LoadInt32(&lc.state) == StoppingState {
			//log.Println("stdAndFileWriter stopped.")
			lc.stopped <- true
		}
	}()

	lock := &sync.Mutex{}

	for {
		select {
		case item, ok := <-lc.items:
			if !ok {
				//log.Println("channel has closed.")
				return
			}

			toStdStr, toFileStr := lc.unpack(item)
			lc.outputConsole(item.level, &toStdStr)
			lc.outputFile(&toFileStr)
			if item.level >= ErrorLevel {
				lc.outputErrFile(&toStdStr)
			}

			if item.level == FatalLevel {
				isFatalLog = true
				return
			}
			lc.put(item)

		case <-flushTicker.C:
			// 定时检查日志文件大小，进行文件切分
			logFileSize := fileSize(lc.pathName)
			logErrFileSize := fileSize(lc.errPathName)
			if logFileSize >= lc.config.singleFileSize || logErrFileSize >= lc.config.singleFileSize {
				lock.Lock()
				if logFileSize >= lc.config.singleFileSize {
					if err := lc.switchFile(); err != nil {
						log.Printf("log switch failed: %v\n", err)
						lock.Unlock()
						panic(err)
					}
				}
				if logErrFileSize >= lc.config.singleFileSize {
					if err := lc.switchErrFile(); err != nil {
						log.Printf("log switch failed: %v\n", err)
						lock.Unlock()
						panic(err)
					}
				}
				lock.Unlock()
			}

		case <-switchTimer.C:
			// 对每天产生的日志文件切分
			lock.Lock()
			if err := lc.switchFile(); err != nil {
				log.Printf("log switch failed: %v\n", err)
				lock.Unlock()
				panic(err)
			}
			if err := lc.switchErrFile(); err != nil {
				log.Printf("log switch failed: %v\n", err)
				lock.Unlock()
				panic(err)
			}
			lock.Unlock()
			switchTimer.Reset(getNextSwitchTime())
		}
	}
}

func (lc *loggerCore) stop() {
	//log.Println("stop")
	//defer log.Println("exit")

	if atomic.LoadInt32(&lc.state) == StoppedState {
		return
	}

	// 设置state为StoppingState，限制畸形儿写入
	atomic.StoreInt32(&lc.state, StoppingState)

	// 关闭通道
	close(lc.items)

	// 等待通道中的数据被读取完毕。缓存通道只有读取完数据才能判断通道已经关闭，即返回ok==false
	<-lc.stopped

	if lc.pFile != nil {
		lc.pFile.Close()
		lc.pFile = nil
	}

	if lc.pErrFile != nil {
		lc.pErrFile.Close()
		lc.pErrFile = nil
	}

	atomic.SwapInt32(&lc.state, StoppedState)
}

// 输出到控制台
func (lc *loggerCore) outputConsole(level logLevel, s *string) {
	stdLogger(level).Println(*s)
}

// 输出到日志文件
func (lc *loggerCore) outputFile(s *string) {
	if lc.pFile != nil {
		fmt.Fprintln(lc.fileWriter, *s)
	}
}

func (lc *loggerCore) outputErrFile(s *string) {
	if lc.pErrFile != nil {
		fmt.Fprintln(lc.errFileWriter, *s)
	}
}

// 解包操作
func (lc *loggerCore) unpack(item *logItem) (toStd string, toFile string) {
	now := time.Now().Format("2006-01-02 15:04:05.000 MST")

	if lc.config.module == defaultNoneModule {
		// [time][level][filename:line]:content
		toStd = fmt.Sprintf("[%s][%s][%s:%d]: %s", now, getLogLevelStr(item.level), item.fileName, item.line, item.content)
		if lc.config.jsonFile {
			toFile = fmt.Sprintf(`{"LEVEL":"%s","Time":"%v","File":"%s","Line":"%d","LocalFunc":"%s","CONTENT":"%s"}`,
				getLogLevelStr(item.level), now, item.fileName, item.line, item.localFunc, item.content)
		} else {
			toFile = toStd
		}
	} else {
		// [time][module][level][filename:line]: content
		toStd = fmt.Sprintf("[%s][%s][%s][%s:%d]: %s", now, lc.config.module, getLogLevelStr(item.level), item.fileName, item.line, item.content)
		if lc.config.jsonFile {
			toFile = fmt.Sprintf(`{"MODULE":"%s","LEVEL":"%s","Time":"%v","File":"%s","Line":"%d","LocalFunc":"%s","CONTENT":"%s"}`,
				lc.config.module, getLogLevelStr(item.level), now, item.fileName, item.line, item.localFunc, item.content)
		} else {
			toFile = toStd
		}
	}

	return
}

// 将item对象放回临时对象池
func (lc *loggerCore) put(l *logItem) {
	lc.pool.Put(l)
}

// 获取item对象
func (lc *loggerCore) get() *logItem {
	return lc.pool.Get().(*logItem)
}

//取基础文件名
func (lc *loggerCore) getFileName(path string) (filename string) {
	now := time.Now()
	filename = path + now.Format("20060102") + "." + lc.localIP //+ ".log"
	return
}

// 切换文件
func (lc *loggerCore) switchFile() error {
	fileName := lc.getFileName(lc.config.path)

	// 确认目录存在
	if err := os.MkdirAll(lc.config.path, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	// 先关闭旧文件再切换
	if lc.pFile != nil {
		if err := lc.pFile.Close(); err != nil {
			return err
		}
		lc.pFile = nil
	}

	// 创建或者打开已存在文件
	file, pathName, err := lc.newFile(fileName)
	if err != nil {
		return err
	}

	lc.pFile = file
	lc.pathName = pathName
	lc.fileWriter = bufio.NewWriter(file)

	return nil
}

// 切换错误记录文件
func (lc *loggerCore) switchErrFile() error {
	fileName := lc.getFileName(lc.config.errPath) + ".error"

	// 确认目录存在
	if err := os.MkdirAll(lc.config.errPath, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	// 先关闭旧文件再切换
	if lc.pErrFile != nil {
		if err := lc.pErrFile.Close(); err != nil {
			return err
		}
		lc.pErrFile = nil
	}

	// 创建或者打开已存在文件
	file, pathName, err := lc.newFile(fileName)
	if err != nil {
		return err
	}

	lc.pErrFile = file
	lc.errPathName = pathName
	lc.errFileWriter = bufio.NewWriter(file)

	return nil
}

// 新建文件，返回描叙符
func (lc *loggerCore) newFile(fileName string) (pFile *os.File, pathName string, err error) {

	pathName = fileName + ".log"

	for fileID := 2; fileExists(pathName); fileID++ {
		pathName = fileName + fmt.Sprintf(".%02d.log", fileID)
	}

	f, err := os.OpenFile(pathName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, "", err
	}

	return f, pathName, nil
}
