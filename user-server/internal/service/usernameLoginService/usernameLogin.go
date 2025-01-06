package usernameLoginService

import (
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg/token"
)

// 用户名密码登录
func UsernameLogin(username string, password string) (
	*pb.UsernameLoginResponse, error) {
	resp := &pb.UsernameLoginResponse{}
	// 根据用户名（唯一）查找，如果用户名不存在就返回用户名不存在，
	// 如果密码为空就返回未设置密码,如果都存在就返回信息
	userId, userName, role, avatarUrl, err := repository.CheckNameAndPwd(username, password)
	if err != nil {
		return nil, err
	}
	//生成双token
	refreshToken, err := token.GenerateRefreshToken(userId, role)
	if err != nil {
		return nil, err
	}
	accessToken, err := token.GenerateAccessToken(userId, role)
	if err != nil {
		return nil, err
	}
	resp.RefreshToken = refreshToken
	resp.AccessToken = accessToken
	//把双token存入redis数据库
	err = repository.SaveToken(userId, resp.RefreshToken, resp.AccessToken)
	if err != nil {
		return nil, err
	}
	return &pb.UsernameLoginResponse{
		UserId:       userId,
		Username:     userName,
		Role:         role,
		Avatar:       avatarUrl,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
