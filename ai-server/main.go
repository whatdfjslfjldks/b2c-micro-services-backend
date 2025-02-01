package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/ai-server/internal/handler"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/ai-server"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50056")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterAIServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50056")

	return grpcServer.Serve(lis)
}

func main() {
	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("ai-server", os.Getenv("api")+":50056", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}
	//localLog.LogLog.Info("etcd: first time register ai-server")

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	//r := gin.Default()
	//r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"*"},                                                                                                        // 允许的跨域来源，* 表示允许所有
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                  // 允许的请求方法
	//	AllowHeaders:     []string{"Access-Token", "Refresh-Token", "X-Real-IP", "X-Forwarded-For", "Origin", "Content-Type", "Authorization"}, // 允许的请求头
	//	AllowCredentials: true,                                                                                                                 // 是否允许携带凭证（例如 cookies）
	//	MaxAge:           12 * time.Hour,                                                                                                       // 预检请求的有效期，单位是时间
	//}))
	//routes.SetupRoutes(r)
	//
	//if err := r.Run(":8080"); err != nil {
	//	log.Fatalf("failed to start HTTP server: %v", err)
	//}

}
