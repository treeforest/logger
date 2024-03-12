package logger

import (
	"testing"
)

func BenchmarkSimpleLogger_Info(b *testing.B) {
	l := NewStdLogger(WithLogLevel(DEBUG), WithPrefix("[BENCHMARK]"), WithShowColor())
	for i := 0; i < b.N; i++ {
		l.Info("Hello, this is benchmark test.")
	}
}

func TestInfo(t *testing.T) {
	l := NewStdLogger(WithPrefix("[TEST]"), WithShowColor())
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")
}
