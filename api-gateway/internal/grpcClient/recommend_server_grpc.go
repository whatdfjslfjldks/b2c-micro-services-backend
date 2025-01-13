package server_grpcClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	recommendServerProto "micro-services/pkg/proto/recommend-server"
)

func (c *GRPCClient) CallRecommendService(serviceName string, method string, request interface{}, response interface{}) error {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}

	//fmt.Println("服务地址-----------------------： ", serviceAddr)

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := recommendServerProto.NewRecommendServiceClient(conn)
	switch method {
	case "clickProduct":
		req := request.(*recommendServerProto.ClickProductRequest)
		_, e := client.ClickProduct(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
	case "browseProduct":
		req := request.(*recommendServerProto.BrowseProductRequest)
		_, e := client.BrowseProduct(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
	case "purchaseProduct":
		req := request.(*recommendServerProto.PurchaseProductRequest)
		_, e := client.PurchaseProduct(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
	case "searchProduct":
		req := request.(*recommendServerProto.SearchProductRequest)
		_, e := client.SearchProduct(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
	default:
		return fmt.Errorf("method %s not supported", method)
	}
	return nil
}
