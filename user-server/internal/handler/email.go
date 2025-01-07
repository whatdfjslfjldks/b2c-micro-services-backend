package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/emailService"
	"micro-services/user-server/internal/service/tokenService"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/token"
)

// 发送邮箱验证码
func (s *Server) EmailSendCode(ctx context.Context, req *pb.EmailSendCodeRequest) (
	*pb.EmailSendCodeResponse, error) {
	//fmt.Println("发送邮箱验证码 入口--------------------")
	email := req.Email
	msg, err, httpCode, statusCode := emailService.SendEmailCode(email)
	if err != nil {
		return &pb.EmailSendCodeResponse{
			Code:       httpCode,
			StatusCode: statusCode,
			Msg:        msg,
		}, nil
	} else {
		return &pb.EmailSendCodeResponse{
			Code:       httpCode,
			StatusCode: statusCode,
			Msg:        msg,
		}, nil
	}
}

// 验证邮箱验证码
func (s *Server) EmailVerifyCode(ctx context.Context, req *pb.EmailVerifyCodeRequest) (
	*pb.EmailVerifyCodeResponse, error) {
	resp := &pb.EmailVerifyCodeResponse{}
	isValid := userPkg.IsEmailValid(req.Email)
	if !isValid {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "邮箱格式错误！"
		return resp, nil
	}
	isVerify := emailService.VerifyCode(req.Email, req.VerifyCode)
	if !isVerify {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "验证码错误或已过期！"
		return resp, nil
	}
	resp, _ = emailService.LoginByEmail(req.Email)
	return resp, nil
}

// 修改邮箱
func (s *Server) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (
	*pb.ChangeEmailResponse, error) {
	resp := &pb.ChangeEmailResponse{}
	isValid := userPkg.IsEmailValid(req.Email)
	if !isValid {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "邮箱格式错误！"
		return resp, nil
	}
	err := emailService.ChangeEmail(req.UserId, req.Email, req.AccessToken)
	if err != nil {
		if err.Error() == "GLB-001" {
			resp.Code = 400
			resp.StatusCode = "GLB-001"
			resp.Msg = "token 已失效或不匹配！"
		} else {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "数据库错误！"
		}
		return resp, nil
	}
	resp.Code = 200
	resp.StatusCode = "GLB-000"
	resp.Msg = "邮箱重置成功！"
	return resp, nil
}

// 通过邮箱重置密码
func (s *Server) ChangePasswordByEmail(ctx context.Context, req *pb.ChangePasswordByEmailRequest) (
	*pb.ChangePasswordByEmailResponse, error) {
	resp := &pb.ChangePasswordByEmailResponse{}

	a := userPkg.IsPasswordValid(req.NewPassword)
	if !a {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "密码不符合要求！"
		return resp, nil
	}
	// 验证token
	_, err := tokenService.TestAccessToken(req.AccessToken)
	if err != nil {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "token 已失效或不匹配!"
		return resp, nil
	}
	claims, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil || claims.UserId != req.UserId {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "token 已失效或不匹配!"
		return resp, nil
	}
	// 验证邮箱验证码
	// TODO 确保在此之前调用了发送验证码接口
	isVerify := emailService.VerifyCode(req.Email, req.VerifyCode)
	if !isVerify {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "验证码错误或已过期！"
		return resp, nil
	}
	// 存储新密码
	err = repository.SaveNewPassword(req.UserId, req.NewPassword)
	if err != nil {
		resp.Code = 500
		resp.StatusCode = "GLB-003"
		resp.Msg = "重置密码失败！"
		return resp, nil
	}
	resp.Code = 200
	resp.StatusCode = "GLB-000"
	resp.Msg = "重置密码成功！"
	return resp, nil
}
