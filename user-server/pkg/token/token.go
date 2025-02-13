package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type MyCustomClaims struct {
	UserId int64  `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TODO 密钥，暂时硬编码，之后用openssl生成rsa密码，私钥存在环境变量里
var secretKey = []byte("secret-key")

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userId int64, role string) (string, error) {
	// 创建 JWT 的有效载荷（claims）

	claims := MyCustomClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 设置当前时间为签发时间
			Issuer:    "b2cPlatform",                  // 签发者
		},
	}
	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名生成 Token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return tokenString, nil
}

// GenerateAccessToken 生成访问令牌，常用于访问受限资源，
func GenerateAccessToken(userId int64, role string) (string, error) {
	// 创建 JWT 的有效载荷（claims）
	claims := MyCustomClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 20)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 设置当前时间为签发时间
			Issuer:    "b2cPlatform",                  // 签发者
		},
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名生成 Token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return tokenString, nil
}

// GetInfoAndCheckExpire 提取 token 里的信息
func GetInfoAndCheckExpire(token string) (
	*MyCustomClaims, error) {
	if token == "" {
		return nil, errors.New("token 不能为空！")
	}
	result, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New("token 解析失败！")
	}
	// 是否过期？
	if claims, ok := result.Claims.(*MyCustomClaims); ok && result.Valid {
		// 已经过期
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token 已过期！") // token 已经过期
		}
		// 未过期，返回基本信息
		return claims, nil
	} else {
		return nil, errors.New("token 不合法！解析失败！")
	}

}
