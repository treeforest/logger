# Logger

Go 轻量级的日志库

## 使用方式

### 导入日志包

```go
import (
	log "github.com/treeforest/logger"
)
```

### 输出到控制台

* 默认控制台输出（同步写）

> ```go
> var defaultLogger Logger = NewStdLogger()
> ...
> log.Debug(...) // std output
> ```

* 用户自定义控制台日志对象

> ```go
> l := log.NewStdLogger(
>   log.WithLevel(log.DEBUG),
>   log.WithPrefix("test"),
> )
> ...
> l.Debug("hello world")
> ```

### 输出到文件  

* 日志分割模式

  > 定时切割：每天24点将会进行日志的切割，把每天的日志分开存储。
  >
  > 阈值切割：在日志写的时候触发检查日志文件大小的事件，若达到阈值，则进行日志切割。

* 日志写模式

  > 同步写：每条日志的打印将会进行一次磁盘的刷新，将缓冲区中的日志刷新到磁盘。
  >
  > 异步写：日志的写入将会先写到缓冲区，只有当异步刷盘的定时器触发时才会将缓冲区的日志刷新到磁盘。使用异步写时，为确保日志不丢失，应使用Stop方法安全关闭日志。

* 同步写日志对象


```go
l := log.NewSyncFileLogger(
		".",
		1024*1024*8,
		log.WithLevel(log.DEBUG),
		log.WithPrefix("example"),
	)
```

* 异步写日志对象

```go
l := log.NewAsyncFileLogger(
		".",
		1024*1024*8,
		1024*64,
		time.Second,
		log.WithLevel(log.DEBUG),
		log.WithPrefix("example"),
	)
```

## 功能特性

- [x] 五种日志级别：debuf/info/warn/error/fatal
- [x] 输出到控制台
- [x] 输出到文件
  - [x] 同步写
  - [x] 异步写
- [x] 日志格式
  - [x] 控制台一般格式
  - [x] json
- [x] 日志文件切割
  - [x] 每天24点切割日志文件
  - [x] 根据日志文件阈值切割
  - [ ] 将ERROR以上级别的日志输出到.error.log文件

## 测试

**测试环境**

> CPU: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz 2.42 GHz
> Memory: 16G
> Go: 1.18.1
> OS: Windows 11
> Hardware: SSD(UMIS RPJTJ512MEE1OWX)

**TPS**

* 日志文件同步写

  ≈ 2000 entry/second

* 日志文件异步写

  16KB缓冲区，每秒进行刷盘：≈ 55000 entry/second

  32KB缓冲区，每秒进行刷盘：≈ 75000 entry/second

  64KB缓冲区，每秒进行刷盘：≈ 110000 entry/second

  128KB缓冲区，每秒进行刷盘：≈ 170000 entry/second

* 控制台输出

  ≈ 10000 entry/second

