package tokenService

import (
	"errors"
	"log"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg/token"
)

// TestAccessToken 验证 accessToken
func TestAccessToken(accessToken string) (bool, error) {
	if accessToken == "" {
		return false, errors.New("token不能为空！")
	}
	// accessToken 不去认证中心（redis）验证
	// 提取 accessToken 的信息，并查看是否过期，过期返回false，没过期看是否在 redis 里
	_, err := token.GetInfoAndCheckExpire(accessToken)
	if err != nil {
		return false, err
	}
	// 查看信息是否与redis存储的匹配
	//err = repository.CheckToken(claims.UserId, accessToken, "accessToken")
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}

// TestRefreshToken 验证 refreshToken
func TestRefreshToken(refreshToken string) (
	string, string, error) {
	if refreshToken == "" {
		return "", "", errors.New("token不能为空！")
	}
	// 提取 refreshToken 的信息，并查看是否过期，过期返回false，没过期看是否在 redis 里
	claims, err := token.GetInfoAndCheckExpire(refreshToken)
	if err != nil {
		log.Println("token.GetInfoAndCheckExpire err: ", err)
		return "", "", err
	}
	// 查看信息是否与redis存储的匹配
	err = repository.CheckToken(claims.UserId, refreshToken, "refreshToken")
	if err != nil {
		log.Println("repository.CheckToken err: ", err)
		return "", "", err
	}
	// 没过期，认证也通过，生成一个 accessToken 返回
	accessToken, err := token.GenerateAccessToken(claims.UserId, claims.Role)
	if err != nil {
		log.Println("token.GenerateAccessToken err: ", err)
		return "", "", err
	}
	rToken, err := token.GenerateRefreshToken(claims.UserId, claims.Role)
	if err != nil {
		log.Println("token.GenerateRefreshToken err: ", err)
		return "", "", err
	}
	//把新的 refreshToken和accessToken 存入redis
	err = repository.SaveToken(claims.UserId, rToken, accessToken)
	if err != nil {
		log.Println("repository.SaveToken err: ", err)
		return "", "", err
	}
	return accessToken, rToken, nil
}
