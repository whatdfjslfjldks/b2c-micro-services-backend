package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	userInfo2 "micro-services/pkg/proto"
	"time"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// 创建 UserServiceClient 客户端
	client := userInfo2.NewUserServiceClient(conn)

	// 准备请求参数
	req := &userInfo2.GetRequest{
		UserId: "123", // 传入需要查询的用户ID
	}

	// 调用 GetUserInfo 方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUserInfo(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetUserInfo: %v", err)
	}

	// 打印返回结果
	fmt.Printf("UserName: %s, UserPwd: %s\n", resp.UserName, resp.UserPwd)
}
