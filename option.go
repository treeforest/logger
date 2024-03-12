package logger

import (
	"fmt"
	"time"
)

const (
	defaultFlushInterval   = time.Second
	defaultRotationTime    = 1    // 1 hour
	defaultRotationSize    = 8    // 8 Mb
	defaultFileBufferBytes = 8192 // 8 KB
)

type LogConfig struct {
	RotationDay     bool          // 是否按天切割日志
	RotationTime    int           // 切割时间（小时）
	RotationSize    int64         // 切割大小（Mb）
	Module          string        // 日志前缀
	LogPath         string        // 日志文件保存路径
	LogLevel        Level         // 日志级别
	FlushInterval   time.Duration // 刷新到磁盘的间隔时间
	ShowColor       bool          // 控制台是否启用颜色输出
	FileBufferBytes int           // 文件写入缓冲区大小（字节）
}

func newConfig() *LogConfig {
	return &LogConfig{
		RotationDay:     false,
		RotationTime:    defaultRotationTime,
		RotationSize:    defaultRotationSize * 1024 * 1024,
		Module:          "",
		LogPath:         ".",
		LogLevel:        DEBUG,
		FlushInterval:   defaultFlushInterval,
		ShowColor:       false,
		FileBufferBytes: defaultFileBufferBytes,
	}
}

type Option func(c *LogConfig)

func WithPrefix(prefix string) Option {
	return func(c *LogConfig) {
		c.Module = fmt.Sprintf("%s ", prefix)
	}
}

func WithLogLevel(level Level) Option {
	return func(c *LogConfig) {
		c.LogLevel = level
	}
}

func WithRotationDay() Option {
	return func(c *LogConfig) {
		c.RotationDay = true
	}
}

func WithRotationTime(rotationTime int) Option {
	return func(c *LogConfig) {
		c.RotationTime = rotationTime
	}
}

func WithRotationSize(rotationSizeInMb int64) Option {
	return func(c *LogConfig) {
		c.RotationSize = rotationSizeInMb * 1024 * 1024
	}
}

func WithFlushInterval(flushInterval time.Duration) Option {
	return func(c *LogConfig) {
		if flushInterval <= time.Millisecond {
			return
		}
		c.FlushInterval = flushInterval
	}
}

func WithShowColor() Option {
	return func(c *LogConfig) {
		c.ShowColor = true
	}
}

func WithFileBufferBytes(fileBufferBytes int) Option {
	return func(c *LogConfig) {
		if fileBufferBytes <= 0 {
			return
		}
		c.FileBufferBytes = fileBufferBytes
	}
}
