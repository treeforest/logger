package log

import (
	"time"
)

const (
	// 管道缓存日志单元的最大容量
	defaultMaxChannelSize uint32 = 1000
	// 检查文件过大的频率
	defaultFileFlushTick time.Duration = 10 * time.Second
	// 默认文件切分大小
	defaultSingleFileSize int64 = 1024 * 1024 * 10
	// 默认无模块显示
	defaultNoneModule string = "NoneModule"
)

type Config struct {
	level          logLevel      // 当前日志打印级别
	path           string        // 日志文件打印路径
	errPath        string        // 错误日志打印路径（默认与path一样）
	singleFileSize int64         // 单个文件大小阈值，达到阈值进行文件切分
	jsonFile       bool          // 是否打印为json格式的日志文件
	fileFlushTick  time.Duration // 检查文件切分的间隔时间
	maxChannelSize uint32        // 管道缓存日志单元的最大容量
	module         string        // 日志打印时输出的模块名
}

func (c *Config) validateConfig() {
	if c.level < DebugLevel || c.level > FatalLevel {
		c.level = DebugLevel
	}

	if c.path != "" {
		rectifyPath(c.path)
		if c.errPath == "" {
			c.errPath = c.path
		}
	}

	if c.errPath != "" {
		rectifyPath(c.errPath)
	}

	if c.singleFileSize < 0 {
		c.singleFileSize = defaultSingleFileSize
	}

	if c.fileFlushTick == time.Duration(0) {
		c.fileFlushTick = defaultFileFlushTick
	}

	if c.maxChannelSize == 0 {
		c.maxChannelSize = defaultMaxChannelSize
	}

	if c.module == "" {
		c.module = defaultNoneModule
	}
}

type Option func(*Config)

func WithFilePath(path string) Option {
	return func(c *Config) {
		c.path = path
	}
}

func WithErrFilePath(path string) Option {
	return func(c *Config) {
		c.errPath = path
	}
}

func WithLogLevel(level logLevel) Option {
	return func(c *Config) {
		c.level = level
	}
}

func WithSingleFileSize(size int64) Option {
	return func(c *Config) {
		c.singleFileSize = size
	}
}

func WithJsonFile(jsonFile bool) Option {
	return func(c *Config) {
		c.jsonFile = jsonFile
	}
}

func WithFileFlushTick(t time.Duration) Option {
	return func(c *Config) {
		c.fileFlushTick = t
	}
}

func WithMaxChannelSize(size uint32) Option {
	return func(c *Config) {
		c.maxChannelSize = size
	}
}

func WithModule(module string) Option {
	return func(c *Config) {
		if module != "" {
			c.module = module
		} else {
			c.module = defaultNoneModule
		}
	}
}