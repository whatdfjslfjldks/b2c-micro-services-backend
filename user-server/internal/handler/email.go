package handler

import (
	"context"
	"database/sql"
	"errors"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/internal/service/emailService"
	"micro-services/user-server/internal/service/tokenService"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/instance"
	"micro-services/user-server/pkg/token"
)

// EmailSendCode 发送邮箱验证码
func (s *Server) EmailSendCode(ctx context.Context, req *pb.EmailSendCodeRequest) (
	*pb.EmailSendCodeResponse, error) {
	for {
		select {
		case <-ctx.Done():
			// 请求已被取消，可能是超时导致的
			return nil, ctx.Err()
		default:
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
	}
}

// EmailVerifyCode 验证邮箱验证码
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
		// TODO 待测试
		// 如果注册过，获取对应id并调用风控，如果没注册过就continue
		id, e := repository.GetUserIdByEmail(req.Email)
		if e == nil {
			r := instance.GrpcClient.RiskIpAndAgentCheck(id, req.Ip, req.UserAgent, "FAIL")
			if r != nil && r.RiskStatus == "RISK" {
				return &pb.EmailVerifyCodeResponse{
					Code:       r.Code,
					StatusCode: r.StatusCode,
					Msg:        r.Msg,
				}, nil
			}
		}
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "验证码错误或已过期！"
		return resp, nil
	}
	resp, _ = emailService.LoginByEmail(req.Email)

	// 风控转发 如果RISK=>拦截
	r := instance.GrpcClient.RiskIpAndAgentCheck(resp.UserId, req.Ip, req.UserAgent, "SUCCESS")
	if r != nil && r.RiskStatus == "RISK" {
		return &pb.EmailVerifyCodeResponse{
			Code:       r.Code,
			StatusCode: r.StatusCode,
			Msg:        r.Msg,
		}, nil
	}
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
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/changeEmail",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
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
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/changePasswordByEmail",
			Source:      "user-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		return resp, nil
	}
	resp.Code = 200
	resp.StatusCode = "GLB-000"
	resp.Msg = "重置密码成功！"
	return resp, nil
}

func (s *Server) GetEmailByUserId(ctx context.Context, req *pb.GetEmailByUserIdRequest) (
	*pb.GetEmailByUserIdResponse, error) {
	resp := &pb.GetEmailByUserIdResponse{}
	email, err := repository.GetEmailByUserId(req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/getEmailByUserId",
			Source:      "user-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		return nil, err
	}
	resp.Email = email
	return resp, nil
}

func (s *Server) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (
	*pb.SendEmailResponse, error) {
	//resp:= &pb.SendEmailResponse{}
	emailService.SendEmail(req.Email, req.Subject, req.Content)
	return nil, nil
}
