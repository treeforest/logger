package logger

import "testing"

func TestFileLogger(t *testing.T) {
	l := NewSyncFileLogger()
	s := "hello world"
	l.Debug(s)
	l.Info(s)
	l.Warn(s)
	l.Error(s)
	// l.Fatal(s)
}
