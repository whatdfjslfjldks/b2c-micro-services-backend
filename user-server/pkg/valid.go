package userPkg

import (
	"github.com/dlclark/regexp2"
	"math/rand"
	"regexp"
	"time"
)

func IsEmailValid(email string) bool {
	reg := `^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	result := regexp.MustCompile(reg)
	return result.MatchString(email)
}

func GenerateVerifyCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("0123456789")
	code := make([]rune, length)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func IsPasswordValid(password string) bool {
	// 密码规则：至少一个字母，一个数字，长度6-15，允许字母、数字、@#等符号
	regex := `^(?=.*[a-zA-Z])(?=.*\d)[a-zA-Z\d@#$%^&*!_+\-=]{6,15}$`
	re := regexp2.MustCompile(regex, 0)

	// 使用正则匹配密码
	match, _ := re.MatchString(password)
	return match
}
