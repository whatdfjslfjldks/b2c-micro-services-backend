package emailService

import (
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
)

func ChangeEmail(userid int64, email string, accessToken string) error {
	// 验证 token
	if ok, err := tokenService.TestAccessToken(accessToken); !ok {
		return err
	}
	// 将邮箱存入数据库
	err := repository.ChangeEmail(userid, email)
	if err != nil {
		return err
	}
	return nil
}
