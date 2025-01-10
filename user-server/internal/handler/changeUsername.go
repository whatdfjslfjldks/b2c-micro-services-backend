package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/service/changeUsernameService"
	"micro-services/user-server/pkg/instance"
)

func (s *Server) ChangeUsername(ctx context.Context, req *pb.ChangeUsernameRequest) (
	*pb.ChangeUsernameResponse, error) {
	resp := &pb.ChangeUsernameResponse{}
	id := req.UserId
	username := utils.Filter(req.Username)
	accessToken := req.AccessToken
	err := changeUsernameService.ChangeUsername(id, username, accessToken)
	if err != nil {
		if err.Error() == "GLB-003" {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "数据库错误！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/changeUsername",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
		} else {
			resp.Code = 400
			resp.StatusCode = "GLB-001"
			resp.Msg = "用户名已存在或 token 失效！"
		}
		return resp, nil
	} else {
		resp.Code = 200
		resp.StatusCode = "GLB-000"
		resp.Msg = "用户名修改成功！"
		resp.Username = username
		return resp, nil
	}
}
