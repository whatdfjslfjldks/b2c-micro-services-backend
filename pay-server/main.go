package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/order-server/pkg/instance"
	"micro-services/pay-server/internal/handler"
	"micro-services/pay-server/internal/notify"
	"micro-services/pay-server/pkg/config"
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

func initConfig() {
	err := config.InitMysqlConfig()
	if err != nil {
		log.Fatalf("Error initializing mysql config: %v", err)
		return
	}
	err = config.InitRedisConfig()
	if err != nil {
		log.Fatalf("Error initializing redis config: %v", err)
		return
	}
	config.InitMySql()
	config.InitRedis()
}

func main() {
	initConfig()
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

	instance.NewInstance()

	// 启动异步通知
	notify.AlipayNotify()

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}

//// TODO TradeCreate 测试
//// 沙盒网关地址 https://openapi-sandbox.dl.alipaydev.com/gateway.do 用不了？ 使用方法是地址拼接返回的交易号参数，到支付宝界面完成支付
//
//// 读取私钥文件
//data, err := ioutil.ReadFile("pay-server/pkg/ali/privateKey.pem")
//if err != nil {
//	log.Fatalf("读取私钥失败: %v", err)
//}
//// 初始化支付宝客户端
//client, err := alipay.New(ali.AppID, string(data), false)
//if err != nil {
//	log.Fatalf("初始化支付宝客户端失败: %v", err)
//}
//// 加载支付宝公钥
//err = client.LoadAliPayPublicKey(ali.AlipayPublicKey)
//if err != nil {
//	log.Fatalf("加载支付宝公钥失败: %v", err)
//}
//
//// 生成一个唯一的订单号
//orderId, err := uuid.NewUUID()
//if err != nil {
//	log.Fatalf("uuid生成失败: %v", err)
//}
//
//r := alipay.TradeCreate{
//	Trade: alipay.Trade{
//		Subject:        "测试 TradeCreate",
//		OutTradeNo:     orderId.String(),
//		TotalAmount:    "0.01",
//		NotifyURL:      "test",
//		ReturnURL:      "http://localhost:3000/",
//		TimeoutExpress: "30m",
//	},
//	BuyerId: ali.UserId,
//}
//resp, err := client.TradeCreate(context.Background(), r)
//
//if resp.Code == "10000" {
//	fmt.Println("支付宝交易号！:", resp.TradeNo)
//} else {
//	fmt.Println("创建失败！")
//}
