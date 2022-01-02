package slog

import (
	"log"
	"os"
)

var l Logger = &simpleLogger{l: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	SetLevel(lvl Level)
}

type Level int

const (
	DEBUG Level = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

var mapping map[Level]string = map[Level]string{
	DEBUG: "DEBU",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERRO",
	FATAL: "FATA",
	PANIC: "PANI",
}

func SetLogger(logger Logger) {
	l = logger
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}
func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	l.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}

func SetLevel(lvl Level) {
	l.SetLevel(lvl)
}