package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/service/changeUsernameService"
)

func (s *Server) ChangeUsername(ctx context.Context, req *pb.ChangeUsernameRequest) (
	*pb.ChangeUsernameResponse, error) {
	resp := &pb.ChangeUsernameResponse{}
	id := req.UserId
	username := utils.Filter(req.Username)
	accessToken := req.AccessToken
	err := changeUsernameService.ChangeUsername(id, username, accessToken)
	if err != nil {
		resp.Msg = err.Error()
		return resp, err
	} else {
		resp.Msg = "修改成功"
		resp.Username = username
		return resp, nil
	}
}
