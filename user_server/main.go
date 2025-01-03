package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/user_server"
	"net"
	"os"
	"time"
)

// server 结构体实现了 UnimplementedUserServiceServer 接口
type server struct {
	pb.UnimplementedUserServiceServer
}

// 发送邮箱验证码
func (s *server) EmailSendCode(ctx context.Context, req *pb.EmailSendCodeRequest) (
	*pb.EmailSendCodeResponse, error) {

	return nil, nil

}

// 创建并启动 gRPC 服务
func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50051")

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
	err = etcdServices.RegisterService("user_server", os.Getenv("api")+":50051")
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
