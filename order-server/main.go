package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/order-server/internal/handler"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/order-server"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50058")

	if err != nil {
		log.Fatalf("forderled to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterOrderServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50058")

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
	err = etcdServices.RegisterService("order-server", os.Getenv("api")+":50058", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("forderled to start gRPC server: %v", err)
	}

}
