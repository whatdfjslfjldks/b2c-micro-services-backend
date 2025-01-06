package changePasswordService

import (
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
)

func ChangePassword(id int64, oldPassword string, newPassword string, accessToken string) error {
	// 先判断 token 是否有效
	_, err := tokenService.TestAccessToken(accessToken)
	if err != nil {
		return err
	}
	// 判断旧密码是否正确
	err = repository.CheckOldPassword(id, oldPassword)
	if err != nil {
		return err
	}
	// 存入新密码
	err = repository.SaveNewPassword(id, newPassword)
	if err != nil {
		return err
	}
	return nil
}
