package handler

import (
	"context"
	"micro-services/order-server/internal/repository"
	pb "micro-services/pkg/proto/order-server"
	"micro-services/user-server/pkg/token"
)

func (s *Server) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (
	*pb.GetOrderDetailResponse, error) {
	// 查验token
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.GetOrderDetailResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.GetOrderDetailResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	// 获取订单项目信息，并查看与id是否符合
	orderDetail, err := repository.GetOrderDetail(req.OrderId, claim.UserId)
	if err != nil {
		return &pb.GetOrderDetailResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "订单不存在或不具有访问权限！" + err.Error(),
		}, nil
	}
	return orderDetail, nil
}
