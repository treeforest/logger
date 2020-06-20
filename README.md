# 日志库

## 使用
~~~
log = logger.NewFileLogger(1 * 1024 * 1024, logger.DebugLevel, "./", "test.log")
defer log.Close()

log.Debug("----这是一条测试的日志----")
log.Info("----这是一条测试的日志----")
log.Warn("----这是一条测试的日志----")
log.Error("----这是一条测试的日志----")
log.Fatal("----这是一条测试的日志----")
~~~