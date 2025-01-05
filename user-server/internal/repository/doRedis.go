package repository

import (
	"errors"
	"fmt"
	"micro-services/user-server/pkg/config"
	"time"
)

func StoreCodeInRedis(email, code string) {
	err := config.RdClient.Set(config.Ctx, email, code, 2*time.Minute).Err() //五分钟过期
	if err != nil {
		fmt.Println("存储验证码失败: ", err)
	} else {
		//fmt.Println("验证码已存储在 Redis 中")
	}
}

func GetCode(email string) (string, error) {
	code, err := config.RdClient.Get(config.Ctx, email).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

func SaveToken(userId int64, refreshToken string, accessToken string) error {
	// 将双token存入redis
	cmd := config.RdClient1.HMSet(config.Ctx,
		fmt.Sprintf("%d", userId),
		"accessToken", accessToken,
		"refreshToken", refreshToken)
	success, err := cmd.Result()
	if err != nil {
		//fmt.Println("Redis error: ", err)
		return errors.New("token存入redis出错")
	}
	if !success {
		//fmt.Println("redissss: ", err)
		return errors.New("token存入redis出错")
	}
	return nil

}
