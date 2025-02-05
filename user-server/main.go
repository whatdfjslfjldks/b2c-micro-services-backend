package main

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/internal/handler"
	"micro-services/user-server/pkg/config"
	"micro-services/user-server/pkg/instance"
	"micro-services/user-server/pkg/localLog"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// 创建并启动 gRPC 服务
func startGRPCServer(port int) error {
	address := ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", port, err)
		return err
	}

	grpcServer := grpc.NewServer()
	// 注册所有服务🌈
	pb.RegisterUserServiceServer(grpcServer, &handler.Server{})

	reflection.Register(grpcServer)
	log.Printf("gRPC server is listening on port %d", port)

	return grpcServer.Serve(lis)
}

func initConfig() {
	err := config.InitEmailConfig()
	if err != nil {
		log.Fatalf("Error initializing internal config: %v", err)
		return
	}
	err = config.InitRedisConfig()
	if err != nil {
		log.Fatalf("Error initializing redis config: %v", err)
		return
	}
	err = config.InitMysqlConfig()
	if err != nil {
		log.Fatalf("Error initializing mysql config: %v", err)
		return
	}
	config.InitRedis()
	config.InitMySql()

	err = localLog.InitLog()
	if err != nil {
		log.Fatalf("Error initializing log config: %v", err)
		return
	}
}

func registerService(etcdServices *etcd.EtcdService, serviceName string, api string, port int, ttl int64) {
	address := api + ":" + strconv.Itoa(port)
	err := etcdServices.RegisterService(serviceName, address, ttl)
	if err != nil {
		log.Fatalf("Error registering service on port %d: %v", port, err)
	}
	l := fmt.Sprintf("etcd: first time register user-server on port: %d ", port)
	localLog.UserLog.Info(l)
}

func main() {
	// 初始化email,redis
	initConfig()
	// 注册服务到 etcd
	etcdServices, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		log.Fatalf("Error creating etcdservice: %v", err)
	}
	defer etcdServices.Close()

	// 要启动的服务实例数量
	instanceCount := 4
	var wg sync.WaitGroup
	wg.Add(instanceCount)

	for i := 0; i < instanceCount; i++ {
		port := 50060 + i
		go func(p int) {
			defer wg.Done()
			id, _ := uuid2.NewUUID()
			// 注册服务到 etcd
			registerService(etcdServices, "user-server"+id.String(), os.Getenv("api"), p, 60)
			instance.NewInstance()
			// 启动 gRPC 服务
			if err := startGRPCServer(p); err != nil {
				log.Fatalf("failed to start gRPC server on port %d: %v", p, err)
			}
		}(port)
	}

	wg.Wait()
}
