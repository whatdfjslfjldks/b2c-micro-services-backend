package handler

import (
	"context"
	"errors"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/service/emailService"
	emailPkg "micro-services/user-server/pkg"
)

// 发送邮箱验证码
func (s *Server) EmailSendCode(ctx context.Context, req *pb.EmailSendCodeRequest) (
	*pb.EmailSendCodeResponse, error) {
	email := req.Email
	msg, err := emailService.SendEmailCode(email)
	if err != nil {
		return &pb.EmailSendCodeResponse{
			Msg: msg,
		}, err
	} else {
		return &pb.EmailSendCodeResponse{
			Msg: msg,
		}, nil
	}
}

// 验证邮箱验证码
func (s *Server) EmailVerifyCode(ctx context.Context, req *pb.EmailVerifyCodeRequest) (
	*pb.EmailVerifyCodeResponse, error) {
	isValid := emailPkg.IsEmailValid(req.Email)
	if !isValid {
		return nil, errors.New("邮箱格式错误！")
	}
	isVerify := emailService.VerifyCode(req.Email, req.VerifyCode)
	if !isVerify {
		return nil, errors.New("验证码错误或已过期！")
	}
	resp, err := emailService.LoginByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
