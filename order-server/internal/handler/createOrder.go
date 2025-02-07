package handler

import (
	"context"
	"fmt"
	pb "micro-services/pkg/proto/order-server"
	"micro-services/user-server/pkg/token"
)

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (
	*pb.CreateOrderResponse, error) {
	// 查验身份，并获取userId
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	fmt.Println("1111: ", req.ProductId)
	fmt.Println("claim:", claim.UserId)

	return nil, nil
}
