package server_grpcClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	productServerProto "micro-services/pkg/proto/product-server"
)

func (c *GRPCClient) CallProductService(serviceName string, method string, request interface{}, response interface{}) error {
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
	client := productServerProto.NewProductServiceClient(conn)
	switch method {
	case "getProductList":
		req := request.(*productServerProto.GetProductListRequest)
		resp, e := client.GetProductList(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.GetProductListResponse) = *resp
	case "uploadProductByExcel":
		req := request.(*productServerProto.UploadProductByExcelRequest)
		resp, e := client.UploadProductByExcel(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.UploadProductByExcelResponse) = *resp
	case "getProductDetailById":
		req := request.(*productServerProto.GetProductDetailByIdRequest)
		resp, e := client.GetProductDetailById(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.GetProductDetailByIdResponse) = *resp
	case "uploadSecKillProduct":
		req := request.(*productServerProto.UploadSecKillProductRequest)
		resp, e := client.UploadSecKillProduct(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.UploadSecKillProductResponse) = *resp
	case "getSecKillList":
		req := request.(*productServerProto.GetSecKillListRequest)
		resp, e := client.GetSecKillList(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.GetSecKillListResponse) = *resp
	case "getProductById":
		req := request.(*productServerProto.GetProductByIdRequest)
		resp, e := client.GetProductById(context.Background(), req)
		if e != nil {
			return fmt.Errorf("failed to call gRPC method: %v", e)
		}
		*response.(*productServerProto.GetProductByIdResponse) = *resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}
	return nil
}
