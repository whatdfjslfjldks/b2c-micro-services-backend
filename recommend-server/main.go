package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/recommend-server/internal/handler"
	"micro-services/recommend-server/pkg/config"
	"micro-services/recommend-server/pkg/instance"
	h "micro-services/recommend-server/pkg/kafka/handler"
	"net"
	"os"
	"time"
)

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50055")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterRecommendServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Println("gRPC server is listening on port 50055")

	return grpcServer.Serve(lis)
}

func initKafka() {
	// 初始化生产者
	h.InitProducer()
}
func initConfig() {
	err := config.InitMysqlConfig()
	if err != nil {
		log.Fatalf("Error initializing internal config: %v", err)
		return
	}
	err = config.InitRedisConfig()
	if err != nil {
		log.Fatalf("Error initializing redis config: %v", err)
		return
	}
	config.InitMySql()
	config.InitRedis()
}

// TODO 基于用户的协同过滤算法
// TODO 存储与当前用户相似用户的userid，数据预热，缓存redis，MySQL做持久化存储，与redis弱一致性
// TODO redis存储用户的商品，权重的map，方便快拿和计算，相似度计算可以做挂起线程，同时控制速度，定期同步redis和mysql
func main() {

	initKafka()
	initConfig()

	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()
	// 注册服务到 etcd
	err = etcdServices.RegisterService("recommend-server", os.Getenv("api")+":50055", 60)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}
	instance.NewInstance()

	// 启动消费者进程
	go handler.ConsumeMsg()

	// 启动 gRPC 服务
	if err := startGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}
