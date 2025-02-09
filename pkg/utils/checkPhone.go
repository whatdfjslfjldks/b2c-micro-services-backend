package utils

import (
	"regexp"
)

// CheckPhone 验证手机号是否符合格式
func CheckPhone(phoneNumber string) bool {
	// 定义正则表达式，匹配中国大陆手机号
	// 以1开头，第二位数字为3-9，后面跟着9位数字
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)

	// 使用正则表达式匹配手机号
	return re.MatchString(phoneNumber)
}
