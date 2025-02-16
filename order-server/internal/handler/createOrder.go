package handler

import (
	"context"
	"log"
	"micro-services/order-server/internal/repository"
	"micro-services/order-server/pkg/instance"
	"micro-services/order-server/pkg/kafka"
	pb "micro-services/pkg/proto/order-server"
	pay "micro-services/pkg/proto/pay-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/pkg/token"
)

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (
	*pb.CreateOrderResponse, error) {
	// 查看手机号是否符合格式
	if !utils.CheckPhone(req.Phone) {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "手机号格式不正确！",
		}, nil
	}
	// 查验身份，并获取userId
	claim, err := token.GetInfoAndCheckExpire(req.AccessToken)
	if err != nil {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "token不可用或已失效！",
		}, nil
	}
	// 查看数据库里有没有用户id
	if !repository.IsUserIdExist(claim.UserId) {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "用户不存在！",
		}, nil
	}
	totalPrice, err := repository.CalcTotalPrice(req.ProductId, req.ProductAmount)
	if err != nil {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        err.Error(),
		}, nil
	}

	// TODO 查看库存够不够并更新库存
	err = repository.CheckProductStock(req.ProductId, req.ProductAmount)
	if err != nil {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        err.Error(),
		}, nil
	}

	// 存储订单基本信息
	orderId, err := repository.CreateOrder(claim.UserId, req.Address, req.Detail, req.Name, req.Phone, req.Note, req.ProductId, req.TypeName, req.ProductAmount, totalPrice)
	if err != nil {
		return &pb.CreateOrderResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        err.Error(),
		}, nil
	}
	request := pay.TradePreCreateRequest{
		Subject:    orderId,
		TotalPrice: totalPrice,
		ReturnUrl:  "test",
		OrderId:    orderId,
	}
	var resp interface{}
	// 调用支付服务，生产alipay支付二维码
	err = instance.GrpcClient.TradePreCreate(&request, &resp)
	if err != nil {
		log.Println("instance.GrpcClient.TradePreCreate err:", err)
		return &pb.CreateOrderResponse{
			Code:       500,
			StatusCode: "GLB-002",
			Msg:        "调用支付服务错误！",
		}, nil
	}
	orderPreCreateResp := resp.(*pay.TradePreCreateResponse)
	if orderPreCreateResp.Code != 200 {
		log.Printf("instance.GrpcClient.TradePreCreate err: %s", orderPreCreateResp.Msg)
		return &pb.CreateOrderResponse{
			Code:       500,
			StatusCode: "PAY-001",
			Msg:        "创建aliPay失败！",
		}, nil
	}
	// 将订单和二维码存入redis，35分钟过期,供用户获取,过期自动删除，包括订单也是过期自动删除
	// 支付二维码是30分钟，redis长5分钟是为了防止订单过期了用户还能支付
	err = repository.SaveAliPayQRCode(orderId, orderPreCreateResp.CodeUrl)
	if err != nil {
		log.Println("repository.SaveAliPayQRCode err:", err)
		return &pb.CreateOrderResponse{
			Code:       500,
			StatusCode: "ORD-001",
			Msg:        "创建订单失败！",
		}, nil
	}

	// 把订单状态发到kafka上去 partition 0 初始状态
	err = kafka.SendMessageToPartition(0, orderId)
	if err != nil {
		log.Println("kafka.SendMessageToPartition err:", err)
		return &pb.CreateOrderResponse{
			Code:       500,
			StatusCode: "ORD-001",
			Msg:        "消息队列状态更新错误！",
		}, nil
	}

	return &pb.CreateOrderResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "创建订单成功",
		OrderId:    orderId,
	}, nil
}
