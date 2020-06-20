# 日志库

## 使用
~~~
log = logger.NewFileLogger(1 * 1024 * 1024, logger.DebugLevel, "./", "test.log")
defer log.Close()

log.Debugf("----这是一条测试的日志----")
log.Infof("----这是一条测试的日志----")
log.Warnf("----这是一条测试的日志----")
log.Errorf("----这是一条测试的日志----")
log.Fatalf("----这是一条测试的日志----")
~~~