package log

import "fmt"

type logger struct {
	core *loggerCore
}

func newLogger(opts ...Option) Logger {
	return &logger{
		core: newLoggerCore(5, opts...),
	}
}

func (l *logger) Debug(a ...interface{}) {
	l.core.debug(a...)
}

func (l *logger) Debugf(format string, a ...interface{}) {
	l.core.debug(fmt.Sprintf(format, a...))
}

func (l *logger) Info(a ...interface{}) {
	l.core.info(a...)
}

func (l *logger) Infof(format string, a ...interface{}) {
	l.core.info(fmt.Sprintf(format, a...))
}

func (l *logger) Warn(a ...interface{}) {
	l.core.warn(a...)
}

func (l *logger) Warnf(format string, a ...interface{}) {
	l.core.warn(fmt.Sprintf(format, a...))
}

func (l *logger) Error(a ...interface{}) {
	l.core.error(a...)
}

func (l *logger) Errorf(format string, a ...interface{}) {
	l.core.error(fmt.Sprintf(format, a...))
}

func (l *logger) Fatal(a ...interface{}) {
	l.core.fatal(a...)
}

func (l *logger) Fatalf(format string, a ...interface{}) {
	l.core.fatal(fmt.Sprintf(format, a...))
}

func (l *logger) SetConfig(opts ...Option) {
	l.core.setConfig(opts...)
}