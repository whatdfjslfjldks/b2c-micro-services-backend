package service

import (
	"micro-services/product-server/pkg/model/enums"
)

// IsSessionValid 判断秒杀场次是否有效
func IsSessionValid(sessionId int32) bool {
	// 检查 SecKillTime 中是否有这个 sessionId
	_, exists := enums.SecKillTime[sessionId]
	return exists
}
