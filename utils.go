package log

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"time"
)

/*
 * 存放公共的工具函数
 */

func getFirstSwitchTime() time.Duration {
	// 到明天凌晨间隔多长时间
	now := time.Now()
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(24 * time.Hour).Sub(now)
}

func getNextSwitchTime() time.Duration {
	return 24 * time.Hour
}

//文件大小
func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 0
	}

	return f.Size()
}

// 文件是否存在，存在返回true
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// 获取本地IP
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	}
	var ip = "localhost"
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

// 目录添加"/"
func rectifyPath(spath string) string {
	if len(spath) >= len("/") && spath[len(spath)-len("/"):] == "/" {
		return spath
	} else {
		return spath + "/"
	}
}

// 获取调用者信息
func getCallerInfo(skip uint32) (fileName, funcName string, line int) {
	pc, file, line, ok := runtime.Caller(int(skip))
	if !ok {
		return
	}

	// 从file(x/y/xx.go)中获取文件名
	fileName = path.Base(file)
	// 根据pc拿到函数名
	funcName = path.Base(runtime.FuncForPC(pc).Name())
	return
}

// 获取日志文件名，以时间为记录节点
func getFileLoggerNameByTime() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	min := now.Minute()
	sec := now.Second()
	var filename string
	filename = fmt.Sprintf("%d", year)
	if month < 10 {
		filename = filename + fmt.Sprintf("0%d", month)
	} else {
		filename = filename + fmt.Sprintf("%d", month)
	}
	if day < 10 {
		filename = filename + fmt.Sprintf("0%d", day)
	} else {
		filename = filename + fmt.Sprintf("%d", day)
	}
	if hour < 10 {
		filename = filename + fmt.Sprintf("0%d", hour)
	} else {
		filename = filename + fmt.Sprintf("%d", hour)
	}
	if min < 10 {
		filename = filename + fmt.Sprintf("0%d", min)
	} else {
		filename = filename + fmt.Sprintf("%d", min)
	}
	if sec < 10 {
		filename = filename + fmt.Sprintf("0%d", sec)
	} else {
		filename = filename + fmt.Sprintf("%d", sec)
	}

	return filename + ".log"
}
