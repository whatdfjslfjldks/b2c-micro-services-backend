package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/handler"
	"micro-services/product-server/pkg/config"
	"micro-services/product-server/pkg/instance"
	"micro-services/product-server/pkg/localLog"
	"net"
	"os"
	"time"
)

// 创建并启动 gRPC 服务
func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterProductServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50054")

	return grpcServer.Serve(lis)
}
func initConfig() {
	err := config.InitMysqlConfig()
	if err != nil {
		log.Fatalf("Error initializing internal config: %v", err)
		return
	}
	config.InitMySql()

	err = localLog.InitLog()
	if err != nil {
		log.Fatalf("Error initializing local log: %v", err)
		return
	}
}
func main() {

	initConfig()

	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}

	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("product-server", os.Getenv("api")+":50054", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	localLog.ProductLog.Info("etcd: first time register product-server")

	instance.NewInstance()

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
