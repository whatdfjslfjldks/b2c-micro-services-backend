package repository

import (
	"fmt"
	"micro-services/recommend-server/pkg/config"
	"micro-services/recommend-server/pkg/kafka/model"
)

func CalAndSaveVectorInToRedis(msg model.Recommend) error {
	fieldName := fmt.Sprintf("product_%d", msg.ProductId)

	// 计算新的 weight 值
	weight := 0
	if msg.Status == "CLICK" {
		weight = 1
	} else if msg.Status == "BROWSE" {
		weight = 2
	} else if msg.Status == "PURCHASE" {
		weight = 3
	}

	// 使用 HIncrBy 实现原有值基础上的加法操作
	cmd := config.RdClient3.HIncrBy(config.Ctx, fmt.Sprintf("%d", msg.UserId), fieldName, int64(weight))

	//key := fmt.Sprintf("%d", msg.UserId)
	//result, _ := config.RdClient3.HGetAll(config.Ctx, key).Result()
	//
	//fmt.Println("sfafqf asfas: ", result)

	return cmd.Err()
}
