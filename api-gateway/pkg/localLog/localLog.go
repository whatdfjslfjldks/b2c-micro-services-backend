package localLog

import (
	"fmt"
	"micro-services/pkg/utils"
)

var GateWayLog *utils.LogWrapper

// InitLog 初始化本地日志
func InitLog() error {
	var err error
	// 初始化日志并赋值给全局变量
	GateWayLog, err = utils.NewLogWrapper("api-gateway")
	if err != nil {
		return fmt.Errorf("初始化日志失败: %v", err)
	}
	return nil
}
