package localLog

import (
	"fmt"
	"micro-services/pkg/utils"
)

var LogLog *utils.LogWrapper

// InitLog 初始化本地日志
func InitLog() error {
	var err error
	// 初始化日志并赋值给全局变量
	LogLog, err = utils.NewLogWrapper("log-server")
	if err != nil {
		return fmt.Errorf("初始化日志失败: %v", err)
	}
	return nil
}
