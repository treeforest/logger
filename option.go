package logger

import "fmt"

type config struct {
	// prefix 日志前缀
	prefix string
	// lvl 日志级别
	lvl Level
}

func newConfig() *config {
	return &config{
		prefix: "",
		lvl:    DEBUG,
	}
}

type Option func(c *config)

func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.prefix = fmt.Sprintf("[%s] ", prefix)
	}
}

func WithLevel(lvl Level) Option {
	return func(c *config) {
		c.lvl = lvl
	}
}
