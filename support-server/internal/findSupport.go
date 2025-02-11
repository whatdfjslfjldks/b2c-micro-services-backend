package internal

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"micro-services/support-server/pkg/config"
)

// FindSupport 查找在线的客服并随机返回一个支持ID
func FindSupport(c *gin.Context) {
	// 获取所有的支持ID（假设有一个合适的 key 格式，例如 "support:online" 存储所有在线客服的 supportID）
	keys, err := config.RdClient.Keys(config.Ctx, "*").Result()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询在线客服失败",
		})
		return
	}

	// 筛选出 status 为 1 的客服
	var onlineSupports []string
	for _, key := range keys {
		status, err := config.RdClient.HGet(config.Ctx, key, "status").Result()
		if err != nil {
			// 跳过没有状态的 key 或错误的 key
			continue
		}
		if status == "1" {
			onlineSupports = append(onlineSupports, key)
		}
	}

	// 如果没有找到在线客服，返回相应提示
	if len(onlineSupports) == 0 {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "没有在线的客服",
		})
		return
	}

	// 随机选择一个在线的客服
	randomIndex := rand.Intn(len(onlineSupports))
	selectedSupportID := onlineSupports[randomIndex]

	// 返回选中的客服ID
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "找到在线客服",
		"data": selectedSupportID,
	})
}
