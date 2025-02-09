package grpcClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"micro-services/pkg/etcd"
	payServerProto "micro-services/pkg/proto/pay-server"
	"time"
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

func (c *GRPCClient) TradePreCreate(request interface{}, response *interface{}) error {
	//获取服务地址
	serviceAddr, err := c.etcdClient.GetService("pay-server")
	if err != nil {
		log.Printf("failed to get service address: %v", err)
	}
	idx := rand.Intn(len(serviceAddr))
	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr[idx], grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		log.Printf("failed to connect to gRPC service: %v\n", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	// 创建 gRPC 客户端
	client := payServerProto.NewPayServiceClient(conn)
	req := request.(*payServerProto.TradePreCreateRequest)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, e := client.TradePreCreate(ctx, req)
	if e != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("请求超时")
		}
		return fmt.Errorf("failed to call sendVerifyCode: %v", err)
	}
	*response = resp
	return nil
}
