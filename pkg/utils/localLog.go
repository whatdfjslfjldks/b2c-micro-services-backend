package utils

import (
	"fmt"
	"log"
	"os"
)

// 存储运行日志到本地

// LogWrapper 封装日志功能
type LogWrapper struct {
	logger *log.Logger
	file   *os.File
}

// NewLogWrapper 创建一个新的日志实例
func NewLogWrapper(serviceName string) (*LogWrapper, error) {
	// 打开日志文件（如果不存在则创建）
	file, err := os.OpenFile("E:/micro-services/"+serviceName+"/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("无法打开日志文件: %v", err)
	}

	// 创建一个新的 logger
	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)

	return &LogWrapper{
		logger: logger,
		file:   file,
	}, nil
}

// Info 输出信息日志
func (lw *LogWrapper) Info(msg string) {
	lw.logger.Println("[INFO] " + msg)
}

// Warn 输出警告日志
func (lw *LogWrapper) Warn(msg string) {
	lw.logger.Println("[WARN] " + msg)
}

// Error 输出错误日志
func (lw *LogWrapper) Error(msg string) {
	lw.logger.Println("[ERROR] " + msg)
}

// Close 关闭文件
func (lw *LogWrapper) Close() {
	lw.file.Close()
}
