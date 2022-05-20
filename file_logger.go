package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// NewAsyncFileLogger 返回一个异步写的日志对象
func NewAsyncFileLogger(path string, capacity int64, bufSize int, flushInterval time.Duration, opts ...Option) Logger {
	l := newFileLogger(path, capacity, bufSize, opts...)
	go l.asyncWrite(flushInterval)
	return l
}

// NewSyncFileLogger 返回一个同步写的日志对象
func NewSyncFileLogger(path string, capacity int64, opts ...Option) Logger {
	l := newFileLogger(path, capacity, 0, opts...)
	go l.syncWrite()
	return l
}

// NewFileLogger 默认日志对象为异步写模式
func NewFileLogger(opts ...Option) Logger {
	return NewAsyncFileLogger(".", 1024*1024*4, 1024*64, time.Second, opts...)
}

func newFileLogger(path string, capacity int64, bufSize int, opts ...Option) *fileLogger {
	pool := &sync.Pool{
		New: func() interface{} {
			return new(logEntry)
		},
	}

	conf := newConfig()
	for _, o := range opts {
		o(conf)
	}

	var c chan *event
	if bufSize > 0 {
		if bufSize > 4096 {
			c = make(chan *event, 4096)
		} else {
			c = make(chan *event, bufSize)
		}
	} else {
		c = make(chan *event)
	}

	l := &fileLogger{
		conf:      conf,
		f:         nil,
		pool:      pool,
		callDepth: 2,
		path:      path,
		filename:  "",
		capacity:  capacity,
		bufSize:   bufSize,
		stop:      make(chan struct{}, 1),
		stopped:   false,
		c:         c, // 不设置缓冲区，禁止异步写
	}
	l.splitLogFile()

	return l
}

type fileLogger struct {
	sync.RWMutex
	conf      *config
	f         *os.File
	bw        *bufio.Writer
	pool      *sync.Pool
	callDepth int
	path      string
	filename  string
	capacity  int64
	bufSize   int
	stop      chan struct{}
	stopOnce  sync.Once
	stopped   bool
	c         chan *event
}

func (l *fileLogger) Debug(v ...interface{}) {
	l.output(DEBUG, fmt.Sprint(v...))
}

func (l *fileLogger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, fmt.Sprintf(format, v...))
}

func (l *fileLogger) Info(v ...interface{}) {
	l.output(INFO, fmt.Sprint(v...))
}

func (l *fileLogger) Infof(format string, v ...interface{}) {
	l.output(INFO, fmt.Sprintf(format, v...))
}

func (l *fileLogger) Warn(v ...interface{}) {
	l.output(WARN, fmt.Sprint(v...))
}

func (l *fileLogger) Warnf(format string, v ...interface{}) {
	l.output(WARN, fmt.Sprintf(format, v...))
}

func (l *fileLogger) Error(v ...interface{}) {
	l.output(ERROR, fmt.Sprint(v...))
}

func (l *fileLogger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
}

func (l *fileLogger) Fatal(v ...interface{}) {
	l.output(FATAL, fmt.Sprint(v...))
	l.Stop()
	os.Exit(1)
}

func (l *fileLogger) Fatalf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
	l.Stop()
	os.Exit(1)
}

func (l *fileLogger) SetLevel(lvl Level) {
	l.Lock()
	l.conf.lvl = lvl
	l.Unlock()
}

type logEntry struct {
	Prefix string `json:"prefix"`
	Time   string `json:"time"`
	Level  string `json:"level"`
	File   string `json:"file"`
	Func   string `json:"func"`
	Msg    string `json:"msg"`
}

type event struct {
	data []byte
	done chan struct{}
}

func (l *fileLogger) isExistFile(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

func (l *fileLogger) output(lvl Level, msg string) {
	if l.stopped {
		return
	}

	l.RLock()
	level := l.conf.lvl
	l.RUnlock()

	if level > lvl {
		return
	}

	var file, fn string = "???", "???"
	var line int = 0
	var pc uintptr
	var ok bool

	pc, file, line, ok = runtime.Caller(l.callDepth)
	if ok {
		file = path.Base(file)
		fn = path.Base(runtime.FuncForPC(pc).Name())
	}

	entry := l.pool.Get().(*logEntry)
	entry.Prefix = l.conf.prefix
	entry.Time = time.Now().Format("2006-01-02 15:04:05.000")
	entry.Level = mapping[lvl]
	entry.File = fmt.Sprintf("%s:%d", file, line)
	entry.Func = fn
	entry.Msg = msg

	b, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	l.pool.Put(entry)

	var e *event
	if l.bufSize > 0 {
		// 异步写
		e = &event{data: b, done: nil}
	} else {
		// 同步写
		e = &event{data: b, done: make(chan struct{}, 1)}
	}

	l.c <- e

	if e.done != nil {
		<-e.done
	}
}

func (l *fileLogger) asyncWrite(flushInterval time.Duration) {
	dayTimer := l.getDayTimer()
	asyncWriteTicker := time.NewTimer(flushInterval) //刷盘间隔

	var err error
	var fi os.FileInfo

	for {
		select {
		case <-l.stop:
			return

		case <-asyncWriteTicker.C:
			_ = l.bw.Flush()
			if err = l.f.Sync(); err != nil {
				panic(err)
			}
			asyncWriteTicker = time.NewTimer(flushInterval)

		case e := <-l.c:
			_, err = l.bw.Write(e.data)
			if err != nil {
				panic(err)
			}
			err = l.bw.WriteByte('\n')
			if err != nil {
				panic(err)
			}

			capacity := atomic.LoadInt64(&l.capacity)
			if capacity <= 0 {
				// 存储无限制
				break
			}
			fi, err = os.Stat(l.filename)
			if err != nil {
				panic(err)
			}
			if fi.Size() > capacity {
				l.Lock()
				l.splitLogFile()
				l.Unlock()
			}

		case <-dayTimer.C:
			time.Sleep(time.Second)
			l.Lock()
			l.splitLogFile()
			l.Unlock()
			dayTimer = l.getDayTimer()
		}
	}
}

func (l *fileLogger) syncWrite() {
	dayTimer := l.getDayTimer()
	var err error
	var fi os.FileInfo

	for {
		select {
		case <-l.stop:
			return

		case e := <-l.c:
			// 1. 输出日志
			_, err = l.bw.Write(e.data)
			if err != nil {
				panic(err)
			}
			err = l.bw.WriteByte('\n')
			if err != nil {
				panic(err)
			}
			// 同步写到磁盘
			_ = l.bw.Flush()
			if err = l.f.Sync(); err != nil {
				panic(err)
			}

			// 通知写成功的事件
			e.done <- struct{}{}

			// 2. 进行日志写时检查：检查文件大小是否达到阈值，若达到阈值，则进行日志文件切割
			capacity := atomic.LoadInt64(&l.capacity)
			if capacity <= 0 {
				// 存储无限制
				break
			}
			fi, err = os.Stat(l.filename)
			if err != nil {
				panic(err)
			}
			if fi.Size() > capacity {
				l.Lock()
				l.splitLogFile()
				l.Unlock()
			}

		case <-dayTimer.C:
			// 1. 每天24点触发，进行日志文件切割
			time.Sleep(time.Second)
			l.Lock()
			l.splitLogFile()
			l.Unlock()
			dayTimer = l.getDayTimer()
		}
	}
}

func (l *fileLogger) getDayTimer() *time.Timer {
	now := time.Now()
	d := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	).Add(24 * time.Hour).Sub(now)
	return time.NewTimer(d)
}

func (l *fileLogger) Stop() {
	l.stopOnce.Do(func() {
		close(l.stop)
		if l.f != nil {
			_ = l.bw.Flush()
			_ = l.f.Close()
		}
		l.stopped = true
	})
}

func (l *fileLogger) splitLogFile() {
	if l.f != nil {
		_ = l.bw.Flush()
		err := l.f.Close()
		if err != nil {
			panic(err)
		}
	}

	filename := l.getFilename()
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	l.filename = filename
	l.f = f
	if l.bufSize > 0 {
		l.bw = bufio.NewWriterSize(f, l.bufSize)
	} else {
		l.bw = bufio.NewWriter(f)
	}
}

func (l *fileLogger) getFilename() string {
	var err error
	var id int

	// 解析出文件的编号
	if l.filename != "" {
		a := strings.Split(l.filename, ".")
		idStr := a[len(a)-2]
		id, err = strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}
	} else {
		id = 1
	}

	var filename string

	// 组装文件名(时间.编号.log)
	prefix := filepath.Join(l.path, time.Now().Format("20060102"))
	for {
		filename = prefix + fmt.Sprintf(".%02d.log", id)
		if !l.isExistFile(filename) {
			break
		}
		id++
	}

	return filename
}
