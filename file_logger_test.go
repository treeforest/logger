package logger

import "testing"

func TestFileLogger(t *testing.T) {
	l := NewSyncFileLogger(".", 1024*1024*4)
	s := "hello world"
	l.Debug(s)
	l.Info(s)
	l.Warn(s)
	l.Error(s)
	// l.Fatal(s)
}
