package repository

import (
	"context"
	"fmt"
	"micro-services/risk-server/pkg/config"
	"strconv"
)

func SaveLoginInfoInToRedis(userId int64, ip string, agent string, status string) error {
	// 构造唯一键（组合userId、ip、agent、status）
	hashKey := fmt.Sprintf("%d:%s:%s:%s", userId, ip, agent, status)
	fmt.Println("生成的hashKey: ", hashKey)

	// 先存储基本信息（login_ip, login_agent, login_status）
	// 使用map传递字段和值
	cmd := config.RdClient2.HMSet(config.Ctx,
		hashKey,
		"login_ip", ip,
		"login_agent", agent,
		"login_status", status,
	)

	// 执行命令并检查是否成功
	if err := cmd.Err(); err != nil {
		fmt.Println("Redis error: ", err)
		return err
	}

	// 使用 HIncrBy 增加 count 字段
	_, err := config.RdClient2.HIncrBy(context.Background(), hashKey, "count", 1).Result()
	if err != nil {
		fmt.Println("Redis error during HINCRBY: ", err)
		return err
	}

	return nil
}

func IsIpAndAgentExists(id int64, ip string, agent string) (bool, error) {
	// 构造唯一的哈希键（user_id:ip:agent）
	hashKey := fmt.Sprintf("%d:%s:%s:%s", id, ip, agent, "SUCCESS")

	cmd := config.RdClient2.HGet(context.Background(), hashKey, "count")
	result, err := cmd.Result()
	if err != nil {
		fmt.Println("是否存在: ", err)
		// 如果是找不到该字段（不存在哈希表），则返回错误
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	count, _ := strconv.Atoi(result)
	fmt.Println("数量数量： ", count)
	if count > 0 {
		return true, nil
	}
	return false, nil
}
