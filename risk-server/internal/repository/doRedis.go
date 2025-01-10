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

	// 先存储基本信息（login_ip, login_agent, login_status）
	cmd := config.RdClient2.HSet(context.Background(),
		hashKey,
		"login_ip", ip,
		"login_agent", agent,
		"login_status", status)

	// 执行命令并检查是否成功
	if err := cmd.Err(); err != nil {
		return err
	}

	// 使用 HINCRBY 增加 count 字段
	_, err := config.RdClient2.HIncrBy(context.Background(), hashKey, "count", 1).Result()
	if err != nil {
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
		// 如果是找不到该字段（不存在哈希表），则返回错误
		return false, err
	}
	count, _ := strconv.Atoi(result)
	if count > 0 {
		return true, nil
	}
	return false, nil
}
