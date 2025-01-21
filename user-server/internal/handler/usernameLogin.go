package handler

import (
	"context"
	"fmt"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/service/usernameLoginService"
)

// 用户名密码登录
func (s *Server) UsernameLogin(ctx context.Context, req *pb.UsernameLoginRequest) (
	*pb.UsernameLoginResponse, error) {
	username := utils.Filter(req.Username)
	password := req.Password
	fmt.Println("asd: ", username, password)
	resp, _ := usernameLoginService.UsernameLogin(username, password, req.Ip, req.UserAgent)

	return resp, nil
}
