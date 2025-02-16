package server_grpcClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	userServerProto "micro-services/pkg/proto/user-server"
	"time"
)

// CallService 调用指定的 gRPC 服务方法
func (c *GRPCClient) CallService(serviceName string, method string, request interface{}, response *interface{}) error {
	//获取服务地址
	serviceAddr, err := c.etcdClient.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service address: %v", err)
	}

	//TODO 负载均衡
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
	client := userServerProto.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 调用服务方法
	switch method {
	case "sendVerifyCode":
		req := request.(*userServerProto.EmailSendCodeRequest)
		resp, err := client.EmailSendCode(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "checkVerifyCode":
		req := request.(*userServerProto.EmailVerifyCodeRequest)
		resp, err := client.EmailVerifyCode(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "loginByPassword":
		req := request.(*userServerProto.UsernameLoginRequest)
		resp, err := client.UsernameLogin(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "testAccessToken":
		req := request.(*userServerProto.TestAccessTokenRequest)
		resp, err := client.TestAccessToken(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "testRefreshToken":
		req := request.(*userServerProto.TestRefreshTokenRequest)
		resp, err := client.TestRefreshToken(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "changeUsername":
		req := request.(*userServerProto.ChangeUsernameRequest)
		resp, err := client.ChangeUsername(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "changeEmail":
		req := request.(*userServerProto.ChangeEmailRequest)
		resp, err := client.ChangeEmail(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "changePassword":
		req := request.(*userServerProto.ChangePasswordRequest)
		resp, err := client.ChangePassword(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "changePasswordByEmail":
		req := request.(*userServerProto.ChangePasswordByEmailRequest)
		resp, err := client.ChangePasswordByEmail(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "editUserInfo":
		req := request.(*userServerProto.EditUserInfoRequest)
		resp, err := client.EditUserInfo(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "getUserInfoByUserId":
		req := request.(*userServerProto.GetUserInfoByUserIdRequest)
		resp, err := client.GetUserInfoByUserId(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "uploadAvatar":
		req := request.(*userServerProto.UploadAvatarRequest)
		resp, err := client.UploadAvatar(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "updateName":
		req := request.(*userServerProto.UpdateNameRequest)
		resp, err := client.UpdateName(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	case "updateBio":
		req := request.(*userServerProto.UpdateBioRequest)
		resp, err := client.UpdateBio(ctx, req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("请求超时")
			}
			return fmt.Errorf("failed to call sendVerifyCode: %v", err)
		}
		*response = resp
	default:
		return fmt.Errorf("method %s not supported", method)
	}

	return nil
}
