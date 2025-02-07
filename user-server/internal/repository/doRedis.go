package repository

import (
	"errors"
	"fmt"
	"micro-services/user-server/pkg/config"
	"time"
)

//	func Test() string {
//		pool := goredis.NewPool(config.RdClient)
//		rs := redsync.New(pool)
//
//		// 获取一个锁对象，指定锁的名称
//		mutex := rs.NewMutex("my-redis-lock")
//
//		// 尝试获取锁，设置锁的过期时间为 10 秒
//		if err := mutex.Lock(); err != nil {
//			panic(err)
//		}
//		fmt.Println("锁锁锁")
//		time.Sleep(2 * time.Second)
//		defer mutex.Unlock()
//		return "yes"
//	}
func StoreCodeInRedis(email, code string) {

	err := config.RdClient.Set(config.Ctx, email, code, 2*time.Minute).Err() //五分钟过期
	if err != nil {
		fmt.Println("存储验证码失败: ", err)
	} else {
		//fmt.Println("验证码已存储在 Redis 中")
	}
}

func GetCode(email string, codeString string) error {
	code, err := config.RdClient.Get(config.Ctx, email).Result()
	if err != nil {
		return err
	}
	if codeString != code {
		return errors.New("验证码错误！")
	}
	// 验证成功后，删除 Redis 中的验证码
	_ = config.RdClient.Del(config.Ctx, email).Err()
	return nil
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

func CheckToken(userid int64, token string, tokenType string) error {
	tokenString, err := config.RdClient1.HGet(config.Ctx, fmt.Sprintf("%d", userid), tokenType).Result()
	if err != nil {
		return err
	}
	if token != tokenString {
		return errors.New("token 不匹配！")
	}
	return nil
}
