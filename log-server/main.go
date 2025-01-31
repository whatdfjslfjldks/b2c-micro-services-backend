package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/log-server/internal/handler"
	"micro-services/log-server/pkg/config"
	h "micro-services/log-server/pkg/kafka/handler"
	"micro-services/log-server/pkg/localLog"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/log-server"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50052")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterLogServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50052")

	return grpcServer.Serve(lis)
}
func initKafka() {
	// 初始化生产者配置，里面已经初始化了kafka配置
	h.InitProducer()
}
func initConfig() {
	if err := config.InitMysqlConfig(); err != nil {
		log.Fatal("Error loading config: ", err)
	}
	config.InitMySql()

	err := localLog.InitLog()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}
func main() {

	// 初始化kafka,生产者和消费者
	initKafka()
	initConfig()

	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("log-server", os.Getenv("api")+":50052", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}
	localLog.LogLog.Info("etcd: first time register log-server")

	// 启动消费者服务
	// 异步执行，防止堵塞进程
	go handler.ConsumeMsg()

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}
