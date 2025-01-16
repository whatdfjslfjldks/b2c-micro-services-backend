package repository

import (
	"encoding/json"
	"fmt"
	"log"
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

func SaveSimUserInToRedis(userId string, simUserId []string) error {
	//fmt.Printf("userId: %s, simUserId: %v\n", userId, simUserId)
	data, err := json.Marshal(simUserId)
	if err != nil {
		log.Printf("Error marshaling data: %v\n", err)
		return err
	}
	cmd := config.RdClient4.Set(config.Ctx, userId, data, 0)
	_, e := cmd.Result()
	if e != nil {
		return cmd.Err()
	}

	//result := config.RdClient4.Get(config.Ctx, userId)
	//
	//// 将结果反序列化为 []string
	//var values []string
	//err = json.Unmarshal([]byte(result.Val()), &values)
	//if err != nil {
	//	log.Fatalf("Error unmarshalling Redis result: %v", err)
	//}
	//
	//// 打印里面的数字
	//for _, val := range values {
	//	log.Println("Extracted number:", val)
	//}
	//
	////r:=result.
	//fmt.Println("result: ", values)

	return nil
}

func GetSimUserId(userId int64) ([]string, error) {
	result := config.RdClient4.Get(config.Ctx, fmt.Sprintf("%d", userId))
	if result.Err() != nil {
		return nil, result.Err()
	}
	var values []string
	err := json.Unmarshal([]byte(result.Val()), &values)
	if err != nil {
		log.Fatalf("Error unmarshalling Redis result: %v", err)
	}
	return values, nil
}

func GetSimProductId(targetUserId int64) (
	[]string, error) {
	result := config.RdClient3.HGetAll(config.Ctx, fmt.Sprintf("%d", targetUserId))
	if result.Err() != nil {
		return nil, result.Err()
	}
	var values []string
	for a, _ := range result.Val() {
		values = append(values, a)
	}
	//fmt.Println("valuesasfsdfdsfa: ", values)
	return values, nil
}
