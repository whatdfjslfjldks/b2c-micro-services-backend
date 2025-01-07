package changeUsernameService

import (
	"errors"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
	"micro-services/user-server/pkg/token"
)

func ChangeUsername(id int64, username string, accessToken string) error {
	//查看用户名是否已经存在
	if repository.IsUsernameExist(username) {
		return errors.New("GLB-001")
	}
	// 查验token
	result, err := tokenService.TestAccessToken(accessToken)
	if err != nil && !result {
		return errors.New("GLB-001")
	}
	claims, err := token.GetInfoAndCheckExpire(accessToken)
	if err != nil || claims.UserId != id {
		return errors.New("GLB-001")
	}
	// 将用户名存入数据库
	err = repository.ChangeUsername(id, username)
	if err != nil {
		return errors.New("GLB-003")
	}
	return nil
}
