package server_grpcClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	productServerProto "micro-services/pkg/proto/product-server"
	"time"
)

func (c *GRPCClient) CallProductService(serviceName string, method string, request interface{}, response *interface{}) error {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}

	// 使用随机算法选择一个服务地址
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
	client := productServerProto.NewProductServiceClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用服务方法
	switch method {
	case "getProductList":
		req := request.(*productServerProto.GetProductListRequest)
		resp, e := client.GetProductList(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call getProductList: %v", e)
		}
		*response = resp
	case "uploadProductByExcel":
		req := request.(*productServerProto.UploadProductByExcelRequest)
		resp, e := client.UploadProductByExcel(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call uploadProductByExcel: %v", e)
		}
		*response = resp
	//case "getProductDetailById":
	//	req := request.(*productServerProto.GetProductDetailByIdRequest)
	//	resp, e := client.GetProductDetailById(ctx, req)
	//	if e != nil {
	//		if errors.Is(e, context.DeadlineExceeded) {
	//			return fmt.Errorf("请求超时")
	//		}
	//		return fmt.Errorf("failed to call getProductDetailById: %v", e)
	//	}
	//	response = resp
	case "uploadSecKillProduct":
		req := request.(*productServerProto.UploadSecKillProductRequest)
		resp, e := client.UploadSecKillProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call uploadSecKillProduct: %v", e)
		}
		*response = resp
	case "getSecKillList":
		req := request.(*productServerProto.GetSecKillListRequest)
		resp, e := client.GetSecKillList(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call getSecKillList: %v", e)
		}
		*response = resp
	case "getProductById":
		req := request.(*productServerProto.GetProductByIdRequest)
		resp, e := client.GetProductById(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call getProductById: %v", e)
		}
		*response = resp
	case "purchaseSecKill":
		req := request.(*productServerProto.PurchaseSecKillRequest)
		resp, e := client.PurchaseSecKill(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call purchaseSecKill: %v", e)
		}
		*response = resp
	case "getOrderConfirmProduct":
		req := request.(*productServerProto.GetOrderConfirmProductRequest)
		resp, e := client.GetOrderConfirmProduct(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call getOrderConfirmProduct: %v", e)
		}
		*response = resp
	case "fuzzySearch":
		req := request.(*productServerProto.FuzzySearchRequest)
		resp, e := client.FuzzySearch(ctx, req)
		if e != nil {
			if errors.Is(e, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call fuzzySearch: %v", e)
		}
		*response = resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}

	return nil
}
