package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/risk-server"
	userServerProto "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/risk-server/internal/repository"
	"micro-services/risk-server/internal/service"
	"micro-services/risk-server/pkg/instance"
)

// RiskLoginIpAndAgent TODO: 如果触发WARN风控，发送邮件提醒，操作地点和时间，触发RISK直接拦截, 标准？
func (s *Server) RiskLoginIpAndAgent(ctx context.Context, req *pb.RiskLoginIpAndAgentRequest) (
	*pb.RiskLoginIpAndAgentResponse, error) {
	// 风控 查找，如何这次的ip和agent在表里没出现过，就warn，邮件提醒 redis
	b, e := repository.IsIpAndAgentExists(req.UserId, req.Ip, req.Agent)
	if e != nil {
		// 数据库错误
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         e.Error(),
			RequestPath: "/riskLoginIpAndAgent",
			Source:      "risk-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		//return nil, nil
	}
	if !b {
		//fmt.Println("之前没有登录过!或登录失败")
		// 风控发warn
		email, err := instance.GrpcClient.GetEmailByUserId(
			&userServerProto.GetEmailByUserIdRequest{
				UserId: req.UserId,
			})
		if err != nil {
			//return nil, nil
		}
		//fmt.Println("邮箱是： ", email)
		service.SendEmail(email, "b2c电商平台 - 异常登录提醒", req.Ip, req.Agent)
	} else {
		//fmt.Println("之前登录过！")
	}

	//fmt.Println("SDfds:", req)
	// 存入mysql，持久化缓存
	err := repository.SaveLoginInfoInToMysql(req.UserId, req.Ip, req.Agent, req.LoginStatus)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/riskLoginIpAndAgent",
			Source:      "risk-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}
	// 存入redis，方便快速获取，避免频繁操作数据库
	err = repository.SaveLoginInfoInToRedis(req.UserId, req.Ip, req.Agent, req.LoginStatus)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/riskLoginIpAndAgent",
			Source:      "risk-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}

	return nil, nil
}
