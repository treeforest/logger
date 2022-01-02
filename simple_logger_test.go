package slog

import (
	"log"
	"os"
	"testing"
)

func BenchmarkSimpleLogger_Info(b *testing.B) {
	ll := &simpleLogger{l: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}
	for i := 0; i < b.N; i++ {
		ll.Info("Hello, this is benchmark test.")
	}
}
