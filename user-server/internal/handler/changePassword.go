package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/service/changePasswordService"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/instance"
)

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (
	*pb.ChangePasswordResponse, error) {
	resp := &pb.ChangePasswordResponse{}
	a := userPkg.IsPasswordValid(req.NewPassword)
	if !a {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "密码不符合要求！"
		return resp, nil
	}
	err := changePasswordService.ChangePassword(req.UserId, req.OldPassword, req.NewPassword, req.AccessToken)
	if err != nil {
		if err.Error() == "GLB-001" {
			resp.Code = 400
			resp.StatusCode = "GLB-001"
			resp.Msg = "密码错误或 token 已失效！"
		} else {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "数据库错误！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/changePassword",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
		}
		return resp, nil
	}
	resp.Code = 200
	resp.StatusCode = "GLB-000"
	resp.Msg = "密码修改成功！"
	return resp, nil
}
