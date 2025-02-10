package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg/token"
)

func (s *Server) GetUserInfoByUserId(ctx context.Context, req *pb.GetUserInfoByUserIdRequest) (
	*pb.GetUserInfoByUserIdResponse, error) {
	// 查验token
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.GetUserInfoByUserIdResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.GetUserInfoByUserIdResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	// 获取用户信息
	avatarUrl, name, email, bio, createAt, err := repository.GetUserInfoByUserId(claim.UserId)
	if err != nil {
		return &pb.GetUserInfoByUserIdResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户信息获取失败！" + err.Error(),
		}, nil
	}
	return &pb.GetUserInfoByUserIdResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "用户信息获取成功！",
		AvatarUrl:  avatarUrl,
		Bio:        bio,
		CreateAt:   createAt,
		Email:      email,
		Name:       name,
		UserId:     claim.UserId,
	}, nil

}
