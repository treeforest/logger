package logger

import "fmt"

type LogConfig struct {
	// RotationDay rotation day
	RotationDay bool
	// RotationTime rotation time Hour
	RotationTime int
	// RotationSize rotation size Mb
	RotationSize int64
	// Module 日志前缀
	Module string
	// LogPath log file save path
	LogPath string
	// LogLevel 日志级别
	LogLevel Level
}

func newConfig() *LogConfig {
	return &LogConfig{
		RotationDay:  true,
		RotationTime: 0,
		RotationSize: 0,
		Module:       "",
		LogPath:      ".",
		LogLevel:     DEBUG,
	}
}

type Option func(c *LogConfig)

func WithPrefix(prefix string) Option {
	return func(c *LogConfig) {
		c.Module = fmt.Sprintf("%s ", prefix)
	}
}

func WithLogLevel(lvl Level) Option {
	return func(c *LogConfig) {
		c.LogLevel = lvl
	}
}

func WithRotationDay(rotationDay bool) Option {
	return func(c *LogConfig) {
		c.RotationDay = rotationDay
	}
}

func WithRotationTime(rotationTime int) Option {
	return func(c *LogConfig) {
		c.RotationTime = rotationTime
	}
}

func WithRotationSize(rotationSize int64) Option {
	return func(c *LogConfig) {
		c.RotationSize = rotationSize * 1024 * 1024
	}
}
