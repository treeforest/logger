# Logger

Go 轻量级日志库

## 功能

* 控制日志打印级别，区分不用环境下的打印
* 异步打印日志
* 定时检测文件大小，若文件大小达到阈值，则拆分日志文件
* 对每日的日志文件进行拆分，便于分辨每天的日志
* ERROR/FATAL 在原来记录的基础上，再输出到.error.log记录文件中
* 控制台不同级别的输出颜色不同
* 默认为控制台输出，若需要文件打印，需调用WithFilePath对输出的文件目录等进行初始化
* 调用log.Stop()优雅退出
* 缓存写机制，保证日志文件的完整性

## 使用
~~~
    defer log.Stop()

	log.SetConfig(
		log.WithLogLevel(log.InfoLevel),
		log.WithFilePath("."))

    log.Debug("Debug Message")
    log.Info("Info Message")
    log.Warn("Warn Message")
    log.Error("Error Message")
    log.Fatal("Fatal Message")
    
    log.Debugf("%s, "Debug Message")
    log.Infof("%s, "Info Message")
    log.Warnf("%s, "Warn Message")
    log.Errorf("%s, "Error Message")
    log.Fatalf("%s, "Fatal Message")
~~~