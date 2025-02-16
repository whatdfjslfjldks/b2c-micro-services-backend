package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg/token"
)

func (s *Server) UpdateBio(ctx context.Context, req *pb.UpdateBioRequest) (
	*pb.UpdateBioResponse, error) {
	// 查验身份，并获取userId
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.UpdateBioResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.UpdateBioResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	err = repository.ChangeBio(claim.UserId, req.Bio)
	if err != nil {
		return &pb.UpdateBioResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "修改失败！",
		}, nil
	}
	return &pb.UpdateBioResponse{
		Code:       200,
		StatusCode: "GLB-001",
		Msg:        "修改成功！",
		Bio:        req.Bio,
	}, nil
}
