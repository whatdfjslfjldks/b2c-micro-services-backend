package server_grpcClient

import (
	"fmt"
	"google.golang.org/grpc"
)

func (c *GRPCClient) CallAIService(serviceName string, method string, request interface{}, response interface{}) error {
	// 获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	//client := aiServerProto.NewAIServiceClient(conn)
	switch method {
	case "talk":
		//req := request.(*aiServerProto.TalkRequest)
		//resp, e := client.Talk(context.Background(), req)
		//if e != nil {
		//	return fmt.Errorf("failed to call gRPC method: %v", e)
		//}
		//*response.(*aiServerProto.TalkResponse) = *resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}
	return nil
}
