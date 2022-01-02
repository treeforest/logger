package slog

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

type simpleLogger struct {
	l   *log.Logger
	lvl Level
}

func (l *simpleLogger) output(lvl Level, msg string) {
	if l.lvl > lvl {
		return
	}

	s := ""
	switch lvl {
	case DEBUG:
		s = l.pack(mapping[DEBUG], msg)
	case INFO:
		s = l.pack(color.GreenString(mapping[INFO]), msg)
	case WARN:
		s = l.pack(color.YellowString(mapping[WARN]), msg)
	case ERROR:
		s = l.pack(color.RedString(mapping[ERROR]), msg)
	case FATAL:
		s = l.pack(color.MagentaString(mapping[FATAL]), msg)
	default:
		fmt.Println("unknown level")
		return
	}

	err := l.l.Output(4, s)
	if err != nil {
		panic(err)
	}
}

func (l *simpleLogger) pack(lvl, msg string) string {
	return fmt.Sprintf("[%s] %s", lvl, msg)
}

func (l *simpleLogger) Debug(v ...interface{}) {
	l.output(DEBUG, fmt.Sprint(v...))
}

func (l *simpleLogger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Info(v ...interface{}) {
	l.output(INFO, fmt.Sprint(v...))
}

func (l *simpleLogger) Infof(format string, v ...interface{}) {
	l.output(INFO, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Warn(v ...interface{}) {
	l.output(WARN, fmt.Sprint(v...))
}

func (l *simpleLogger) Warnf(format string, v ...interface{}) {
	l.output(WARN, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Error(v ...interface{}) {
	l.output(ERROR, fmt.Sprint(v...))
}

func (l *simpleLogger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Fatal(v ...interface{}) {
	l.output(FATAL, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *simpleLogger) Fatalf(format string, v ...interface{}) {
	l.output(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *simpleLogger) SetLevel(lvl Level) {
	l.lvl = lvl
}
