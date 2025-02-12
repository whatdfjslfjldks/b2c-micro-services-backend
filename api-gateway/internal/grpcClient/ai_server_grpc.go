package server_grpcClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	aiServerProto "micro-services/pkg/proto/ai-server"
	"time"
)

func (c *GRPCClient) CallAIService(serviceName string, method string, request interface{}, response interface{}) error {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}

	// 使用随机算法选择一个服务地址
	idx := rand.Intn(len(serviceAddr))

	//与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr[idx], grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 创建 gRPC 客户端
	client := aiServerProto.NewAIServiceClient(conn)
	switch method {
	case "talk":
		req := request.(*aiServerProto.TalkRequest)
		resp, e := client.Talk(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("failed to call gRPC method: %v", e)
			}
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		response = resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}
	return nil
}
