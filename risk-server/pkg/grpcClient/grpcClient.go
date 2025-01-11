package grpcClient

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"micro-services/pkg/etcd"
	logServerProto "micro-services/pkg/proto/log-server"
	userServerProto "micro-services/pkg/proto/user-server"
)

type GRPCClient struct {
	etcdClient *etcd.EtcdService
}

// NewGRPCClient 构造 GRPCClient 实例
func NewGRPCClient(etcdClient *etcd.EtcdService) *GRPCClient {
	return &GRPCClient{
		etcdClient: etcdClient,
	}
}

func (c *GRPCClient) PostLog(request interface{}) {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService("log-server")
	if err != nil {
		log.Printf("failed to get service address: %v", err)
	}

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		log.Printf("failed to connect to gRPC service: %v\n", err)
		return
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	client := logServerProto.NewLogServiceClient(conn)
	req := request.(*logServerProto.PostLogRequest)
	_, _ = client.PostLog(context.Background(), req)
}
func (c *GRPCClient) GetEmailByUserId(request interface{}) (
	string, error) {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService("user-server")
	if err != nil {
		log.Printf("failed to get service address: %v", err)
	}

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		log.Printf("failed to connect to gRPC service: %v\n", err)
		return "", err
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	client := userServerProto.NewUserServiceClient(conn)
	req := request.(*userServerProto.GetEmailByUserIdRequest)
	r, e := client.GetEmailByUserId(context.Background(), req)
	if e != nil {
		log.Printf("failed to connect to gRPC service: %v\n", e)
		return "", err
	}
	return r.Email, nil
}

func (c *GRPCClient) SendEmail(request interface{}) {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService("user-server")
	if err != nil {
		log.Printf("failed to get service address: %v", err)
	}

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		log.Printf("failed to connect to gRPC service: %v\n", err)
		return
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	client := userServerProto.NewUserServiceClient(conn)
	req := request.(*userServerProto.SendEmailRequest)
	_, _ = client.SendEmail(context.Background(), req)
}
