package handler

import (
	"context"
	"errors"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/emailService"
	"micro-services/user-server/internal/service/tokenService"
	emailPkg "micro-services/user-server/pkg"
)

// 发送邮箱验证码
func (s *Server) EmailSendCode(ctx context.Context, req *pb.EmailSendCodeRequest) (
	*pb.EmailSendCodeResponse, error) {
	//fmt.Println("发送邮箱验证码 入口--------------------")
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

// 修改邮箱
func (s *Server) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (
	*pb.ChangeEmailResponse, error) {
	resp := &pb.ChangeEmailResponse{}
	isValid := emailPkg.IsEmailValid(req.Email)
	if !isValid {
		resp.Msg = "邮箱格式错误！"
		return resp, errors.New("邮箱格式错误！")
	}
	err := emailService.ChangeEmail(req.UserId, req.Email, req.AccessToken)
	if err != nil {
		resp.Msg = err.Error()
		return resp, err
	}
	resp.Msg = "邮箱重置成功！"
	return resp, nil
}

// 通过邮箱重置密码
func (s *Server) ChangePasswordByEmail(ctx context.Context, req *pb.ChangePasswordByEmailRequest) (
	*pb.ChangePasswordByEmailResponse, error) {
	resp := &pb.ChangePasswordByEmailResponse{}
	// 验证token
	_, err := tokenService.TestAccessToken(req.AccessToken)
	if err != nil {
		resp.Msg = "token 验证失败!"
		return resp, err
	}
	// 验证邮箱验证码
	// TODO 确保在此之前调用了发送验证码接口
	isVerify := emailService.VerifyCode(req.Email, req.VerifyCode)
	if !isVerify {
		resp.Msg = "验证码错误或已过期！"
		return resp, errors.New("验证码错误或已过期！")
	}
	// 存储新密码
	err = repository.SaveNewPassword(req.UserId, req.NewPassword)
	if err != nil {
		resp.Msg = "重置密码失败！"
		return resp, err
	}
	resp.Msg = "重置密码成功！"
	return resp, nil
}
