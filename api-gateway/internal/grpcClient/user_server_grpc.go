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

	fmt.Println("服务地址-----------------------： ", serviceAddr)

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := userServerProto.NewUserServiceClient(conn)

	// TODO ctx可以设置超时时间，但是对应方法那边需要做对应处理，这里先不做处理，得加钱
	// 调用服务方法
	// TODO 写法太冗余，待优化
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
	case "testAccessToken":
		req := request.(*userServerProto.TestAccessTokenRequest)
		resp, err := client.TestAccessToken(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call testAccessToken:%v", err)
		}
		*response.(*userServerProto.TestAccessTokenResponse) = *resp
	case "testRefreshToken":
		req := request.(*userServerProto.TestRefreshTokenRequest)
		resp, err := client.TestRefreshToken(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call testRefreshToken:%v", err)
		}
		*response.(*userServerProto.TestRefreshTokenResponse) = *resp
	case "changeUsername":
		req := request.(*userServerProto.ChangeUsernameRequest)
		resp, err := client.ChangeUsername(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call changeUsername:%v", err)
		}
		*response.(*userServerProto.ChangeUsernameResponse) = *resp
	case "changeEmail":
		req := request.(*userServerProto.ChangeEmailRequest)
		resp, err := client.ChangeEmail(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call changeEmail:%v", err)
		}
		*response.(*userServerProto.ChangeEmailResponse) = *resp
	case "changePassword":
		req := request.(*userServerProto.ChangePasswordRequest)
		resp, err := client.ChangePassword(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call changePassword:%v", err)
		}
		*response.(*userServerProto.ChangePasswordResponse) = *resp
	case "changePasswordByEmail":
		req := request.(*userServerProto.ChangePasswordByEmailRequest)
		resp, err := client.ChangePasswordByEmail(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call changePasswordByEmail:%v", err)
		}
		*response.(*userServerProto.ChangePasswordByEmailResponse) = *resp
	case "editUserInfo":
		req := request.(*userServerProto.EditUserInfoRequest)
		resp, err := client.EditUserInfo(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to call editUserInfo:%v", err)
		}
		*response.(*userServerProto.EditUserInfoResponse) = *resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}

	return nil
}
