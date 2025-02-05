package server_grpcClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	recommendServerProto "micro-services/pkg/proto/recommend-server"
	"time"
)

func (c *GRPCClient) CallRecommendService(serviceName string, method string, request interface{}, response *interface{}) error {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}
	idx := rand.Intn(len(serviceAddr))
	// 与 gRPC 服务建立连接
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

	// 创建 gRPC 客户端
	client := recommendServerProto.NewRecommendServiceClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用服务方法
	switch method {
	case "clickProduct":
		req := request.(*recommendServerProto.ClickProductRequest)
		_, e := client.ClickProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call clickProduct: %v", e)
		}
	case "browseProduct":
		req := request.(*recommendServerProto.BrowseProductRequest)
		_, e := client.BrowseProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call browseProduct: %v", e)
		}
	case "purchaseProduct":
		req := request.(*recommendServerProto.PurchaseProductRequest)
		_, e := client.PurchaseProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call purchaseProduct: %v", e)
		}
	case "searchProduct":
		req := request.(*recommendServerProto.SearchProductRequest)
		_, e := client.SearchProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call searchProduct: %v", e)
		}
	case "getRecommendProductList":
		req := request.(*recommendServerProto.GetRecommendProductListRequest)
		resp, e := client.GetRecommendProductList(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call getRecommendProductList: %v", e)
		}
		*response = resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}

	return nil
}
