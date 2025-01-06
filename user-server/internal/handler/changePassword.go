package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/service/changePasswordService"
)

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (
	*pb.ChangePasswordResponse, error) {
	resp := &pb.ChangePasswordResponse{}
	err := changePasswordService.ChangePassword(req.UserId, req.OldPassword, req.NewPassword, req.AccessToken)
	if err != nil {
		resp.Msg = "密码修改失败！"
		return resp, err
	}
	resp.Msg = "密码修改成功！"
	return resp, nil
}
