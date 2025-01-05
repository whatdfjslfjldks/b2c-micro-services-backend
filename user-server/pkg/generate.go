package userPkg

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成一个随机字母字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" // 随机字母的字符集
	rand.Seed(time.Now().UnixNano())             // 设置随机数种子

	// 构建随机字符串
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))] // 从 charset 中随机选择一个字符
	}
	return string(result)
}

// 生成唯一的用户名
func GenerateUsername() string {
	// 随机字母的长度可以调整，这里是生成 3 个字母
	randomStr := generateRandomString(3)
	username := "用户" + randomStr
	fmt.Println("生成的用户名: ", username)
	return username
}
