package changePasswordService

import (
	"errors"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/tokenService"
	"micro-services/user-server/pkg/token"
)

func ChangePassword(id int64, oldPassword string, newPassword string, accessToken string) error {
	// 先判断 token 是否有效
	_, err := tokenService.TestAccessToken(accessToken)
	if err != nil {
		return errors.New("GLB-001")
	}
	claims, err := token.GetInfoAndCheckExpire(accessToken)
	if err != nil || claims.UserId != id {
		return errors.New("GLB-001")
	}
	// 判断旧密码是否正确
	err = repository.CheckOldPassword(id, oldPassword)
	if err != nil {
		return errors.New("GLB-001")
	}
	// 存入新密码
	err = repository.SaveNewPassword(id, newPassword)
	if err != nil {
		return errors.New("GLB-003")
	}
	return nil
}
