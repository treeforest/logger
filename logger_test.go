package logger

import (
	"testing"
	)

func TestConstLevel(t *testing.T) {

	t.Logf("%v", DebugLevel)
	t.Logf("%v", InfoLevel)
	t.Logf("%v", WarnLevel)
}