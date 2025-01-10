package utils

import "time"

func GetTime() string {
	// 获取当前时间，可以统一配置时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}
