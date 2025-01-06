package changeUsernameService

import (
	"errors"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
)

func ChangeUsername(id int64, username string, accessToken string) error {
	//查看用户名是否已经存在
	if repository.IsUsernameExist(username) {
		return errors.New("用户名已被占用！")
	}
	// 查验token
	result, err := tokenService.TestAccessToken(accessToken)
	if err != nil && !result {
		return err
	}
	// 将用户名存入数据库
	err = repository.ChangeUsername(id, username)
	if err != nil {
		return err
	}
	return nil
}
