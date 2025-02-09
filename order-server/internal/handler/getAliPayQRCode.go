package handler

import (
	"context"
	"log"
	"micro-services/order-server/internal/repository"
	pb "micro-services/pkg/proto/order-server"
)

func (s *Server) GetAliPayQRCode(ctx context.Context, req *pb.GetAliPayQRCodeRequest) (
	*pb.GetAliPayQRCodeResponse, error) {
	// 查找redis，获取对应的阿里支付二维码
	codeUrl, err := repository.GetAliPayQRCode(req.OrderId)
	if err != nil {
		log.Println("repository.GetAliPayQRCode err:", err)
		return &pb.GetAliPayQRCodeResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "订单已过期或不存在！",
		}, nil
	}
	return &pb.GetAliPayQRCodeResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "获取成功！",
		CodeUrl:    codeUrl,
	}, nil
}
