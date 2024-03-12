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
	"time"
)

// NewAsyncFileLogger 返回一个异步写的日志对象
func NewAsyncFileLogger(opts ...Option) Logger {
	l := newFileLogger(true, opts...)
	go l.asyncWrite()
	return l
}

// NewSyncFileLogger 返回一个同步写的日志对象
func NewSyncFileLogger(opts ...Option) Logger {
	l := newFileLogger(false, opts...)
	go l.syncWrite()
	return l
}

// NewFileLogger 默认日志对象为异步写模式
func NewFileLogger(opts ...Option) Logger {
	return NewAsyncFileLogger(opts...)
}

func newFileLogger(async bool, opts ...Option) *fileLogger {
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
	if async {
		c = make(chan *event, 8192)
	} else {
		c = make(chan *event)
	}

	l := &fileLogger{
		conf:      conf,
		f:         nil,
		pool:      pool,
		callDepth: 2,
		filename:  "",
		stop:      make(chan struct{}, 1),
		routineWG: sync.WaitGroup{},
		c:         c, // 不设置缓冲区，禁止异步写
	}
	l.splitLogFile()

	return l
}

type fileLogger struct {
	conf      *LogConfig
	f         *os.File
	bw        *bufio.Writer
	pool      *sync.Pool
	callDepth int
	filename  string
	stop      chan struct{}
	routineWG sync.WaitGroup
	stopOnce  sync.Once
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
	l.output(FATAL, fmt.Sprintf(format, v...))
	l.Stop()
	os.Exit(1)
}

func (l *fileLogger) SetLevel(lvl Level) {
	l.conf.LogLevel = lvl
}

type logEntry struct {
	Module string `json:"module,omitempty"`
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
	select {
	case <-l.stop:
		return
	default:
	}

	if l.conf.LogLevel > lvl {
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
	entry.Module = l.conf.Module
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
	if cap(l.c) > 0 {
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

func (l *fileLogger) asyncWrite() {
	l.routineWG.Add(1)
	defer l.routineWG.Done()

	dayTimer := l.getDayTimer()
	hourTimer := l.getHourTimer()
	flushTimer := time.NewTimer(l.conf.FlushInterval) //刷盘间隔

	var err error
	var fi os.FileInfo

	for {
		select {
		case <-l.stop:
			return

		case <-flushTimer.C:
			_ = l.bw.Flush()
			if err = l.f.Sync(); err != nil {
				panic(err)
			}
			flushTimer = time.NewTimer(l.conf.FlushInterval)

		case e := <-l.c:
			_, err = l.bw.Write(e.data)
			if err != nil {
				panic(err)
			}
			err = l.bw.WriteByte('\n')
			if err != nil {
				panic(err)
			}

			capacity := l.conf.RotationSize
			if capacity <= 0 {
				// 存储无限制
				break
			}
			fi, err = os.Stat(l.filename)
			if err != nil {
				panic(err)
			}
			if fi.Size() > capacity {
				l.splitLogFile()
			}

		case <-dayTimer.C:
			l.splitLogFile()
			dayTimer = l.getDayTimer()

		case <-hourTimer.C:
			l.splitLogFile()
			hourTimer = l.getHourTimer()
		}
	}
}

func (l *fileLogger) syncWrite() {
	l.routineWG.Add(1)
	defer l.routineWG.Done()

	dayTimer := l.getDayTimer()
	hourTimer := l.getHourTimer()

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
			capacity := l.conf.RotationSize
			if capacity <= 0 {
				// 存储无限制
				break
			}

			fi, err = os.Stat(l.filename)
			if err != nil {
				panic(err)
			}
			if fi.Size() > capacity {
				l.splitLogFile()
			}

		case <-dayTimer.C:
			l.splitLogFile()
			dayTimer = l.getDayTimer()

		case <-hourTimer.C:
			l.splitLogFile()
			hourTimer = l.getHourTimer()
		}
	}
}

func (l *fileLogger) getHourTimer() *time.Timer {
	if l.conf.RotationTime == 0 {
		return l.getNeverTriggerTimer()
	}

	now := time.Now()
	d := time.Date(
		now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location(),
	).Add(time.Hour * time.Duration(l.conf.RotationTime)).Sub(now)
	return time.NewTimer(d)
}

func (l *fileLogger) getDayTimer() *time.Timer {
	if !l.conf.RotationDay {
		return l.getNeverTriggerTimer()
	}

	now := time.Now()
	d := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	).Add(24 * time.Hour).Sub(now)
	return time.NewTimer(d)
}

func (l *fileLogger) getFlushTimer() *time.Timer {
	return time.NewTimer(l.conf.FlushInterval)
}

func (l *fileLogger) getNeverTriggerTimer() *time.Timer {
	timer := time.NewTimer(time.Duration(0))
	if !timer.Stop() {
		<-timer.C
	}
	return timer
}

func (l *fileLogger) Stop() {
	l.stopOnce.Do(func() {
		close(l.stop)

		// 等待协程退出
		l.routineWG.Wait()

		if l.f == nil {
			return
		}

		close(l.c)

		// 读出缓冲区的所有日志，并写入文件
		for e := range l.c {
			_, err := l.bw.Write(e.data)
			if err != nil {
				panic(err)
			}
			err = l.bw.WriteByte('\n')
			if err != nil {
				panic(err)
			}
			if e.done != nil {
				e.done <- struct{}{}
			}
		}

		_ = l.bw.Flush()
		_ = l.f.Sync()
		_ = l.f.Close()
		l.bw = nil
		l.f = nil
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
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l.filename = filename
	l.f = f
	if cap(l.c) > 0 {
		l.bw = bufio.NewWriterSize(f, l.conf.FileBufferBytes)
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
	prefix := filepath.Join(l.conf.LogPath, time.Now().Format("2006010215"))
	for {
		filename = prefix + fmt.Sprintf(".%02d.log", id)
		if !l.isExistFile(filename) {
			break
		} else {
			peekFilename := prefix + fmt.Sprintf(".%02d.log", id+1)
			if !l.isExistFile(peekFilename) {
				capacity := l.conf.RotationSize
				if capacity <= 0 {
					break
				}
				fi, err := os.Stat(filename)
				if err != nil {
					panic(err)
				}
				if fi.Size() <= capacity {
					// 日志未达到切割大小，可以继续追加
					break
				}
			}
		}
		id++
	}

	return filename
}
