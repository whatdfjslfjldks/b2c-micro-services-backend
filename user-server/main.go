package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/handler"
	"micro-services/user-server/pkg/config"

	"net"
	"os"
	"time"
)

// 创建并启动 gRPC 服务
func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册发送邮箱验证码
	pb.RegisterUserServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50051")

	return grpcServer.Serve(lis)
}
func initConfig() {
	err := config.InitEmailConfig()
	if err != nil {
		log.Fatalf("Error initializing internal config: %v", err)
		return
	}
	err = config.InitRedisConfig()
	if err != nil {
		log.Fatalf("Error initializing redis config: %v", err)
		return
	}
	err = config.InitMysqlConfig()
	if err != nil {
		log.Fatalf("Error initializing mysql config: %v", err)
		return
	}
	config.InitRedis()
	config.InitMySql()
}

func main() {
	// 初始化email,redis
	initConfig()
	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()

	// 注册服务到 etcd
	err = etcdServices.RegisterService("user-server", os.Getenv("api")+":50051")
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
