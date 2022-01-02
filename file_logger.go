package slog

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileLogger struct {
	lvl      Level
	f        *os.File
	w        *bufio.Writer
	pool     *sync.Pool
	mu       sync.Mutex
	cancel   context.CancelFunc
	path     string
	filename string
	capacity int64
}

func (l *FileLogger) Debug(v ...interface{}) {
	l.output(DEBUG, fmt.Sprint(v...))
}

func (l *FileLogger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, fmt.Sprintf(format, v...))
}

func (l *FileLogger) Info(v ...interface{}) {
	l.output(INFO, fmt.Sprint(v...))
}

func (l *FileLogger) Infof(format string, v ...interface{}) {
	l.output(INFO, fmt.Sprintf(format, v...))
}

func (l *FileLogger) Warn(v ...interface{}) {
	l.output(WARN, fmt.Sprint(v...))
}

func (l *FileLogger) Warnf(format string, v ...interface{}) {
	l.output(WARN, fmt.Sprintf(format, v...))
}

func (l *FileLogger) Error(v ...interface{}) {
	l.output(ERROR, fmt.Sprint(v...))
}

func (l *FileLogger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
}

func (l *FileLogger) Fatal(v ...interface{}) {
	l.output(FATAL, fmt.Sprint(v...))
	l.stop()
	os.Exit(1)
}

func (l *FileLogger) Fatalf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
	l.stop()
	os.Exit(1)
}

func (l *FileLogger) SetLevel(lvl Level) {
	l.lvl = lvl
}

type logItem struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	File  string `json:"file"`
	Func  string `json:"func"`
	Msg   string `json:"msg"`
}

func NewFileLogger(path string, capacity int64) *FileLogger {
	pool := &sync.Pool{
		New: func() interface{} {
			return new(logItem)
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	l := &FileLogger{
		lvl:      DEBUG,
		f:        nil,
		w:        nil,
		pool:     pool,
		mu:       sync.Mutex{},
		cancel:   cancel,
		path:     path,
		filename: "",
		capacity: capacity, // 5kb
	}
	l.newLogFile()
	go l.splitEveryDay(ctx)
	return l
}

func (l *FileLogger) splitEveryDay(ctx context.Context) {
	now := time.Now()
	d := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	).Add(24 * time.Hour).Sub(now)
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
			t.Reset(24 * time.Hour)
			time.Sleep(time.Second)
			l.mu.Lock()
			l.newLogFile()
			l.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (l *FileLogger) newLogFile() {
	if l.f != nil {
		err := l.f.Close()
		if err != nil {
			panic(err)
		}
	}

	prefix := filepath.Join(l.path, time.Now().Format("20060102"))

	filename := ""
	id := 1

	for {
		filename = prefix + fmt.Sprintf(".%02d.log", id)
		if !l.isExistFile(filename) {
			break
		}
		id++
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	l.filename = filename
	l.f = f
	l.w = bufio.NewWriter(f)
}

func (l *FileLogger) isExistFile(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

func (l *FileLogger) output(lvl Level, msg string) {
	if l.lvl > lvl {
		return
	}

	defer func() {
		err := l.w.Flush()
		if err != nil {
			panic(err)
		}
	}()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 如果文件大小达到阈值，则对日志文件进行切割
	fi, err := os.Stat(l.filename)
	if err != nil {
		panic(err)
	}
	if fi.Size() > l.capacity {
		l.newLogFile()
	}

	var file, fun string = "???", "???"
	var line int = 0

	pc, file, line, ok := runtime.Caller(3)
	if ok {
		file = path.Base(file)
		fun = path.Base(runtime.FuncForPC(pc).Name())
	}

	o := l.pool.Get().(*logItem)
	defer l.pool.Put(o)

	o.Time = time.Now().Format("2006-01-02 15:04:05.000")
	o.Level = mapping[lvl]
	o.File = fmt.Sprintf("%s:%d", file, line)
	o.Func = fun
	o.Msg = msg

	buf, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	_, err = l.w.Write(append(buf, '\n'))
	if err != nil {
		panic(err)
	}
}

func (l *FileLogger) stop() {
	_ = l.f.Close()
	l.cancel()
	time.Sleep(time.Millisecond * 50)
}
