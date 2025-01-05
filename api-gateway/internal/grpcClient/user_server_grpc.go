package user_server_grpcClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"micro-services/pkg/etcd"
	userServerProto "micro-services/pkg/proto/user-server"
)

// GRPCClient 封装了与 gRPC 服务的连接
type GRPCClient struct {
	etcdClient *etcd.EtcdService
}

// NewGRPCClient 构造 GRPCClient 实例
func NewGRPCClient(etcdClient *etcd.EtcdService) *GRPCClient {
	return &GRPCClient{
		etcdClient: etcdClient,
	}
}

// CallService 调用指定的 gRPC 服务方法
func (c *GRPCClient) CallService(serviceName, method string, request interface{}, response interface{}) error {
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
	client := userServerProto.NewUserServiceClient(conn)

	// 调用服务方法
	switch method {
	case "sendVerifyCode":
		req := request.(*userServerProto.EmailSendCodeRequest)
		resp, err := client.EmailSendCode(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response.(*userServerProto.EmailSendCodeResponse) = *resp
	case "checkVerifyCode":
		req := request.(*userServerProto.EmailVerifyCodeRequest)
		resp, err := client.EmailVerifyCode(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call checkVerifyCode: %v", err)
		}
		*response.(*userServerProto.EmailVerifyCodeResponse) = *resp
	case "loginByPassword":
		req := request.(*userServerProto.UsernameLoginRequest)
		resp, err := client.UsernameLogin(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call checkVerifyCode: %v", err)
		}
		*response.(*userServerProto.UsernameLoginResponse) = *resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}

	return nil
}
