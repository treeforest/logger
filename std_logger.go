package logger

import (
	"fmt"
	"log"
	"os"
)

func NewStdLogger(opts ...Option) Logger {
	conf := newConfig()
	for _, o := range opts {
		o(conf)
	}
	return &stdLogger{
		l:         log.New(os.Stderr, conf.Module, log.LstdFlags|log.Lshortfile),
		c:         conf,
		callDepth: 3,
	}
}

type stdLogger struct {
	l         *log.Logger
	c         *LogConfig
	callDepth int
}

func (l *stdLogger) output(lvl Level, msg string) {
	if l.c.LogLevel > lvl {
		return
	}

	s := ""
	if l.c.ShowColor {
		s = l.pack(colorMapping[lvl], msg)
	} else {
		s = l.pack(mapping[lvl], msg)
	}

	err := l.l.Output(l.callDepth, s)
	if err != nil {
		panic(err)
	}
}

func (l *stdLogger) pack(lvl, msg string) string {
	return fmt.Sprintf("[%s] %s", lvl, msg)
}

func (l *stdLogger) Debug(v ...interface{}) {
	l.output(DEBUG, fmt.Sprint(v...))
}

func (l *stdLogger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Info(v ...interface{}) {
	l.output(INFO, fmt.Sprint(v...))
}

func (l *stdLogger) Infof(format string, v ...interface{}) {
	l.output(INFO, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Warn(v ...interface{}) {
	l.output(WARN, fmt.Sprint(v...))
}

func (l *stdLogger) Warnf(format string, v ...interface{}) {
	l.output(WARN, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Error(v ...interface{}) {
	l.output(ERROR, fmt.Sprint(v...))
}

func (l *stdLogger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Fatal(v ...interface{}) {
	l.output(FATAL, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *stdLogger) Fatalf(format string, v ...interface{}) {
	l.output(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *stdLogger) SetLevel(lvl Level) {
	l.c.LogLevel = lvl
}

func (l *stdLogger) Stop() {
	// do nothing
}
