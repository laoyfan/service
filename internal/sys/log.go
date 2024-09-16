package sys

import (
	"fmt"
	"log"
	"os"
	"service/util"
	"sync"
)

var (
	logger *log.Logger
	file   *os.File
	once   sync.Once
)

func init() {
	once.Do(func() {
		if ok, _ := util.PathExists("log"); !ok {
			fmt.Println("创建日志文件夹", "log")
			err := os.Mkdir("log", os.ModePerm)
			if err != nil {
				fmt.Println("创建日志文件夹失败", err)
			}
		}
		var err error
		// 打开一个文件用于写入日志
		file, err = os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("无法打开日志文件")
			//log.Fatalf("无法打开日志文件: %v", err)
		}
		logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	})
}

func Log(level string, msg string, keyAndValue ...interface{}) {
	// 构建日志信息
	logMsg := msg
	for i := 0; i < len(keyAndValue); i += 2 {
		logMsg += " " + keyAndValue[i].(string) + "=" + keyAndValue[i+1].(string)
	}
	// 根据日志级别记录日志
	switch level {
	case "DEBUG":
		logger.Println("DEBUG: " + logMsg)
	case "INFO":
		logger.Println("INFO: " + logMsg)
	case "WARN":
		logger.Println("WARN: " + logMsg)
	case "ERROR":
		logger.Println("ERROR: " + logMsg)
	default:
		logger.Println(logMsg)
	}
}

func CloseLog() {
	// 接收到信号后关闭日志文件
	err := file.Close()
	if err != nil {
		fmt.Println("关闭文件出错")
		//log.Printf("关闭文件出错: %v", err)
	}
}
