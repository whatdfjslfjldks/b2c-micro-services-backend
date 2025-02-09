package handler

import (
	"context"
	"fmt"
	"log"
	"micro-services/order-server/internal/repository"
	"micro-services/order-server/pkg/kafka"
	pb "micro-services/pkg/proto/order-server"
)

func (s *Server) TestPaySuccess(ctx context.Context, req *pb.TestPaySuccessRequest) (
	*pb.TestPaySuccessResponse, error) {

	fmt.Printf("收到支付成功消息: %v", req)
	// 消费掉未支付（初始状态下）Kafka的消息
	err := kafka.ConsumePartition(0, req.OrderId)
	if err != nil {
		log.Printf("消费消息失败: %v", err)
		return &pb.TestPaySuccessResponse{
			Code:       500,
			Msg:        "消费消息失败" + err.Error(),
			StatusCode: "ORD-001",
		}, nil
	}
	fmt.Printf("1")
	// 发送新消息到已支付状态的Kafka
	err = kafka.SendMessageToPartition(1, req.OrderId)
	if err != nil {
		log.Printf("发送消息到Kafka失败: %v", err)
		return &pb.TestPaySuccessResponse{
			Code:       500,
			Msg:        "发送消息到消息队列失败",
			StatusCode: "ORD-001",
		}, nil
	}
	fmt.Printf("11")
	// 修改数据库中订单状态
	err = repository.ReverseOrderStatus(req.OrderId, 1)
	if err != nil {
		log.Println("修改数据库订单状态失败")
		return &pb.TestPaySuccessResponse{
			Code:       500,
			Msg:        "修改数据库订单状态失败" + err.Error(),
			StatusCode: "GLB-003",
		}, nil
	}
	fmt.Printf("111")
	// 删除redis里的订单二维码
	err = repository.DeleteAliPayQRCode(req.OrderId)
	if err != nil {
		log.Println("删除redis订单二维码失败")
		return &pb.TestPaySuccessResponse{
			Code:       500,
			Msg:        "修改数据库订单状态失败" + err.Error(),
			StatusCode: "GLB-003",
		}, nil
	}
	fmt.Printf("1111")
	returnUrl := fmt.Sprintf("/orderDetail?orderId=%s", req.OrderId)
	return &pb.TestPaySuccessResponse{
		Code:       200,
		Msg:        "success",
		ReturnUrl:  returnUrl,
		StatusCode: "200",
	}, nil
}
