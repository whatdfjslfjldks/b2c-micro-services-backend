package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pay-server/internal/handler"
	"micro-services/pay-server/internal/notify"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/pay-server"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50057")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterPayServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50057")

	return grpcServer.Serve(lis)
}

func main() {
	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("pay-server", os.Getenv("api")+":50057", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}
	//localLog.LogLog.Info("etcd: first time register ai-server")

	// 启动异步通知
	notify.AlipayNotify()

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}

//func main() {
//	// 读取私钥文件
//	data, err := ioutil.ReadFile("pay-server/pkg/ali/privateKey.pem")
//	if err != nil {
//		log.Fatalf("读取私钥失败: %v", err)
//	}
//
//	// 初始化支付宝客户端
//	client, err := alipay.New(ali.AppID, string(data), false)
//	if err != nil {
//		log.Fatalf("初始化支付宝客户端失败: %v", err)
//	}
//
//	// 加载支付宝公钥
//	err = client.LoadAliPayPublicKey(ali.AlipayPublicKey)
//	if err != nil {
//		log.Fatalf("加载支付宝公钥失败: %v", err)
//	}
//
//	// 生成一个唯一的订单号
//	uuid, err := uuid.NewUUID()
//	if err != nil {
//		log.Fatalf("uuid生成失败: %v", err)
//	}
//
//	//req := alipay.TradeCreate{
//	//	Trade: alipay.Trade{
//	//		Subject:        "测试",
//	//		OutTradeNo:     uuid.String(),
//	//		TotalAmount:    "0.01",
//	//		NotifyURL:      "https://7ea5-103-151-173-97.ngrok-free.app/notify",
//	//		ReturnURL:      "http://localhost:3000/",
//	//		TimeoutExpress: "30m",
//	//	},
//	//	BuyerOpenId: "123",
//	//}
//	// 创建统一收单线下交易预创建请求
//	req := alipay.TradePreCreate{
//		Trade: alipay.Trade{
//			Subject:        "测试",
//			OutTradeNo:     uuid.String(),
//			TotalAmount:    "0.01",
//			NotifyURL:      "https://7ea5-103-151-173-97.ngrok-free.app/notify",
//			ReturnURL:      "http://localhost:3000/",
//			TimeoutExpress: "30m",
//		},
//	}
//
//	ctx := context.Background()
//
//	//resp, err := client.TradeCreate(ctx, req)
//	// 发起预创建交易请求
//	resp, err := client.TradePreCreate(ctx, req)
//	if err != nil {
//		log.Fatalf("发起预创建交易请求失败: %v", err)
//	}
//
//	fmt.Println("resp: ", resp)
//
//	// 处理响应结果
//	//if resp.Code == "10000" {
//	//	// 返回交易二维码链接
//	//	fmt.Printf("交易预创建成功，二维码链接: %s\n", resp.QRCode)
//	//} else {
//	//	fmt.Printf("交易预创建失败，错误码: %s，错误信息: %s\n", resp.Code, resp.Msg)
//	//}
//
//	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
//		// 获取异步通知参数
//		r.ParseForm()
//		noti, err := client.DecodeNotification(r.Form)
//		if err != nil {
//			log.Printf("解析支付宝通知失败: %v", err)
//			http.Error(w, "解析通知失败", http.StatusBadRequest)
//			return
//		}
//		// 验证签名
//		err = client.VerifySign(r.Form)
//		if err != nil {
//			log.Printf("支付宝通知签名验证失败: %v", err)
//			http.Error(w, "验证失败", http.StatusBadRequest)
//			return
//		}
//
//		// 处理不同的交易状态
//		switch noti.TradeStatus {
//		case "WAIT_BUYER_PAY":
//			log.Printf("交易创建，等待买家付款，订单号: %s\n", noti.OutTradeNo)
//		case "TRADE_SUCCESS":
//			log.Printf("支付成功，订单号: %s\n", noti.OutTradeNo)
//			// 在这里可以进行支付成功的业务处理，例如更新数据库订单状态
//		case "TRADE_CLOSED":
//			log.Printf("交易关闭，订单号: %s\n", noti.OutTradeNo)
//		default:
//			log.Printf("未知交易状态: %s，订单号: %s\n", noti.TradeStatus, noti.OutTradeNo)
//		}
//
//		// 响应支付宝，告知已经收到通知
//		w.Write([]byte("success"))
//	})
//
//	// 启动 HTTP 服务
//	go func() {
//		if err := http.ListenAndServe(":8080", nil); err != nil {
//			log.Fatalf("HTTP 服务启动失败: %v", err)
//		}
//	}()
//
//	// 等待异步通知进行支付结果处理
//	select {} // 防止主线程退出
//}
