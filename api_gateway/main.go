package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api_gateway/internal/routes"
)

func main() {
	// 启动 HTTP 服务
	r := gin.Default()
	// 用gin做网关，进行路由的接收和转发
	routes.SetupRoutes(r)

	//// 定义一个路由来处理前端请求
	//r.GET("/v1/userinfo", func(c *gin.Context) {
	//	// 从查询参数获取 user_id
	//	userId := c.DefaultQuery("user_id", "")
	//	if userId == "" {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
	//		return
	//	}
	//
	//	// 创建上下文
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//
	//	// 从 etcd 获取服务地址
	//	etcdServices, err := etcd.NewEtcdService([]string{"localhost:2379"}, 5*time.Second)
	//	if err != nil {
	//		log.Printf("Error creating etcd service: %v", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	//		return
	//	}
	//	defer etcdServices.Close()
	//
	//	// 获取注册的服务地址
	//	serviceAddress, err := etcdServices.GetService("user_server")
	//	if err != nil {
	//		log.Printf("Error fetching service address from etcd: %v", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	//		return
	//	}
	//
	//	// 连接到 gRPC 服务
	//	conn, err := grpc.DialContext(ctx, serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	//	if err != nil {
	//		log.Printf("Failed to dial gRPC server: %v", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	//		return
	//	}
	//	defer conn.Close()
	//
	//	// 创建 gRPC 客户端
	//	client := userInfo2.NewUserServiceClient(conn)
	//
	//	// 调用 GetUserInfo 方法
	//	req := &userInfo2.GetRequest{
	//		UserId: userId, // 使用查询参数中的 user_id
	//	}
	//	resp, err := client.GetUserInfo(ctx, req)
	//	if err != nil {
	//		log.Printf("Failed to call GetUserInfo: %v", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	//		return
	//	}
	//
	//	// 返回从 gRPC 服务获取的用户名和密码
	//	c.JSON(http.StatusOK, gin.H{
	//		"user_name": resp.UserName,
	//		"user_pwd":  resp.UserPwd,
	//	})
	//})

	// 启动 HTTP 服务，监听 8080 端口
	log.Println("HTTP server is listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
