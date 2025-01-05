package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/service/usernameLoginService"
)

// 用户名密码登录
func (s *Server) UsernameLogin(ctx context.Context, req *pb.UsernameLoginRequest) (
	*pb.UsernameLoginResponse, error) {
	username := req.Username
	password := req.Password
	resp, err := usernameLoginService.UsernameLogin(username, password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
