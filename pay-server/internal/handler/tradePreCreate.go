package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/smartwalle/alipay/v3"
	"io/ioutil"
	"log"
	"micro-services/pay-server/pkg/ali"
	pb "micro-services/pkg/proto/pay-server"
	"strconv"
)

func (s *Server) TradePreCreate(ctx context.Context, req *pb.TradePreCreateRequest) (
	*pb.TradePreCreateResponse, error) {
	// 读取私钥文件
	data, err := ioutil.ReadFile("pay-server/pkg/ali/privateKey.pem")
	if err != nil {
		log.Fatalf("读取私钥失败: %v", err)
	}
	// 初始化支付宝客户端
	client, err := alipay.New(ali.AppID, string(data), false)
	if err != nil {
		log.Fatalf("初始化支付宝客户端失败: %v", err)
	}
	// 加载支付宝公钥
	err = client.LoadAliPayPublicKey(ali.AlipayPublicKey)
	if err != nil {
		log.Fatalf("加载支付宝公钥失败: %v", err)
	}

	// 生成一个唯一的订单号
	orderId, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("uuid生成失败: %v", err)
	}

	str := strconv.FormatFloat(float64(req.TotalAmount), 'f', 2, 32)
	// 创建统一收单线下交易预创建请求
	r := alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:        req.Subject,
			OutTradeNo:     orderId.String(),
			TotalAmount:    str,
			NotifyURL:      "https://ff20-103-151-173-97.ngrok-free.app/notify",
			ReturnURL:      "http://localhost:3000/",
			TimeoutExpress: "30m",
		},
	}

	c := context.Background()

	//resp, err := client.TradeCreate(ctx, req)
	// 发起预创建交易请求
	resp, err := client.TradePreCreate(c, r)
	if err != nil {
		log.Fatalf("发起预创建交易请求失败: %v", err)
	}

	// 处理响应结果
	if resp.Code == "10000" {
		// 返回交易二维码链接
		fmt.Printf("交易预创建成功，二维码链接: %s\n", resp.QRCode)
		return &pb.TradePreCreateResponse{
			Code:       200,
			CodeUrl:    resp.QRCode,
			Msg:        "交易预创建成功",
			StatusCode: "GLB-000",
		}, nil
	} else {
		fmt.Printf("交易预创建失败，错误码: %s，错误信息: %s\n", resp.Code, resp.Msg)
		return &pb.TradePreCreateResponse{
			Code:       500,
			CodeUrl:    "",
			Msg:        resp.Msg,
			StatusCode: "PAY-001",
		}, nil
	}

}
