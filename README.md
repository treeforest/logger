# Logger
Go 轻量级的日志库

## 使用方式

1. 在代码里导入

```go
import (
	log "github.com/treeforest/logger"
)
```

2. 输出到控制台

```go
package main

import (
	log "github.com/treeforest/logger"
)

func main() {
	log.Info("Hello World!")
}
```

3. 输出到文件

```go
package main

import (
	log "github.com/treeforest/logger"
)

func main() {   
	log.SetLogger(log.NewFileLogger(".", 1024*1024*5))
	log.Info("Hello World!")
}
```



## 功能特性

- [x] 日志级别控制

- [x] 打印控制台日志

  - [x] 不同级别打印不同的颜色
- [x] 打印json文件日志
  - [x] 定时（每日凌晨）切割日志文件，便于查看每日的日志
  - [x] 根据文件大小阈值切割文件
  - [ ] 将错误级别的日志输出到error文件

