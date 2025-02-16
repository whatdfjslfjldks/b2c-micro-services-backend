package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg/token"
)

func (s *Server) UpdateName(ctx context.Context, req *pb.UpdateNameRequest) (
	*pb.UpdateNameResponse, error) {
	// 查验身份，并获取userId
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.UpdateNameResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.UpdateNameResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	if repository.IsUsernameExist(req.Name) {
		return &pb.UpdateNameResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户名已存在！",
		}, nil
	}
	err = repository.ChangeUsername(claim.UserId, req.Name)
	if err != nil {
		return &pb.UpdateNameResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "修改失败！",
		}, nil
	}
	return &pb.UpdateNameResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "修改成功！",
		Name:       req.Name,
	}, nil
}
