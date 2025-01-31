package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/risk-server"
	"micro-services/risk-server/internal/handler"
	"micro-services/risk-server/pkg/config"
	"micro-services/risk-server/pkg/instance"
	"micro-services/risk-server/pkg/localLog"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50053")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterRiskServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50053")

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
	err = localLog.InitLog()
	if err != nil {
		log.Fatalf("Error initializing log config: %v", err)
		return
	}
}
func main() {
	// TODO 设计风控模块数据库，存用户的ip，常用设备（user-agent），短时间修改密码次数（这个是举个例子）等信息
	// TODO 同时加载常查询数据进入redis，方便快速获取

	// 初始化mysql,redis
	initConfig()

	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("risk-server", os.Getenv("api")+":50053", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	localLog.RiskLog.Info("etcd: first time register risk-server")

	instance.NewInstance()
	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
