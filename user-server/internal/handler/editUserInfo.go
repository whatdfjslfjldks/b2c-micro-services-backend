package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/service/editUserInfo"
	"micro-services/user-server/internal/service/tokenService"
	"micro-services/user-server/pkg/token"
)

func (s *Server) EditUserInfo(ctx context.Context, req *pb.EditUserInfoRequest) (
	*pb.EditUserInfoResponse, error) {
	resp := &pb.EditUserInfoResponse{}
	// 查验token
	if ok, _ := tokenService.TestAccessToken(req.AccessToken); !ok {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "token 过期或不匹配！"
		return resp, nil
	}
	claims, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil || claims.UserId != req.UserId {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "token 过期或不匹配！"
		return resp, nil
	}
	// 过滤 bio
	bio := utils.Filter(req.Bio)
	a := editUserInfo.EditUserInfo(req.UserId, req.AvatarUrl, bio, req.Location)
	return a, nil
}
