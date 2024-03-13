package logger

import (
	"testing"
)

func TestFileLogger(t *testing.T) {
	l := NewSyncFileLogger()
	s := "hello world"
	l.Debug(s)
	l.Info(s)
	l.Warn(s)
	l.Error(s)
	// l.Fatal(s)
}

func BenchmarkAsyncFileLogger(b *testing.B) {
	l := NewAsyncFileLogger(WithLogPath("./tmp"), WithRotationSize(64))
	for i := 0; i < b.N; i++ {
		l.Info("hello world")
	}
	l.Stop()
}

func BenchmarkSyncFileLogger(b *testing.B) {
	l := NewSyncFileLogger()
	for i := 0; i < b.N; i++ {
		l.Info("hello world")
	}
	l.Stop()
}
