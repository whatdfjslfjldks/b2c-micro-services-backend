package utils

import (
	"fmt"
	"time"
)

func GetTime() string {
	// 获取当前时间，可以统一配置时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func ConvertTime(timeStr string) (string, error) {
	// 解析 ISO 8601 格式的时间
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "", err
	}
	// 转换为 MySQL 可接受的格式 YYYY-MM-DD HH:MM:SS
	result := parsedTime.Format("2006-01-02 15:04:05")
	return result, nil
}
