package emailService

import (
	"errors"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
	"micro-services/user-server/pkg/token"
)

func ChangeEmail(userid int64, email string, accessToken string) error {
	// 验证 token
	if ok, _ := tokenService.TestAccessToken(accessToken); !ok {
		return errors.New("GLB-001")
	}

	claims, err := token.GetInfoAndCheckExpire(accessToken)
	if err != nil || claims.UserId != userid {
		return errors.New("GLB-001")
	}
	// 将邮箱存入数据库
	err = repository.ChangeEmail(userid, email)
	if err != nil {
		return errors.New("GLB-003")
	}
	return nil
}
