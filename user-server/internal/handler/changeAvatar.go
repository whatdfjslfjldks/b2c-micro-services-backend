package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/minio"
	"micro-services/user-server/pkg/token"
	"strconv"
)

func (s *Server) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (
	*pb.UploadAvatarResponse, error) {
	// 查验身份，并获取userId
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.UploadAvatarResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.UploadAvatarResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	result, msg, path := minio.UploadFileToMinio(req.File, strconv.FormatInt(claim.UserId, 10)+"touxiang", "/avatar/", "b2c")
	if !result {
		return &pb.UploadAvatarResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        msg,
		}, nil
	}
	p := "b2c" + path
	// 头像存入数据
	err = repository.ChangeAvatar(claim.UserId, p)
	if err != nil {
		return &pb.UploadAvatarResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "头像上传失败！",
		}, nil
	}
	return &pb.UploadAvatarResponse{
		Code:       200,
		StatusCode: "GLB-001",
		Msg:        "上传成功！",
		AvatarUrl:  p,
	}, nil

}
