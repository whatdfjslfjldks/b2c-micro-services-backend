package handler

import (
	"context"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/service/tokenService"
)

// 验证 accessToken 是否过期

func (s *Server) TestAccessToken(ctx context.Context, req *pb.TestAccessTokenRequest) (
	*pb.TestAccessTokenResponse, error) {
	resp := &pb.TestAccessTokenResponse{}
	result, err := tokenService.TestAccessToken(req.AccessToken)
	if err != nil && !result {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "accessToken 已过期或不匹配！"
		return resp, nil
	} else {
		resp.Code = 200
		resp.StatusCode = "GLB-000"
		resp.Msg = "token 验证成功"
		return resp, nil
	}
}

// 验证 refreshToken 是否过期
func (s *Server) TestRefreshToken(ctx context.Context, req *pb.TestRefreshTokenRequest) (
	*pb.TestRefreshTokenResponse, error) {
	resp := &pb.TestRefreshTokenResponse{}
	accessToken, err := tokenService.TestRefreshToken(req.RefreshToken)
	if err != nil {
		resp.Code = 400
		resp.StatusCode = "GLB-001"
		resp.Msg = "refreshToken 已过期或不匹配！"
		return resp, nil
	} else {
		resp.Code = 200
		resp.StatusCode = "GLB-000"
		resp.Msg = "token 验证成功"
		resp.AccessToken = accessToken
		return resp, nil
	}
}
