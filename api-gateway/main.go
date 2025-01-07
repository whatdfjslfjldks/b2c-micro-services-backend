package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api-gateway/internal/instance"
	"micro-services/api-gateway/internal/routes"
)

func main() {
	// 初始化 kafka-user-server
	//err := pkgConfig.InitKafkaConfig()
	//if err != nil {
	//	return
	//}
	//// 发布消息
	//err = pkgConfig.PublishMessage("Hello Kafka!")
	//if err != nil {
	//	log.Fatalf("Error publishing message: %v", err)
	//}
	//
	//// 订阅并消费消息
	//err = pkgConfig.ConsumeMessages()
	//if err != nil {
	//	log.Fatalf("Error consuming messages: %v", err)
	//}
	// 启动 HTTP 服务
	r := gin.Default()
	//创建实例
	instance.NewInstance()
	// 用gin做网关，进行路由的接收和转发
	routes.SetupRoutes(r)

	// 启动 HTTP 服务，监听 8080 端口
	log.Println("HTTP server is listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	// TODO 风控模块，日志模块，订单模块，支付模块
	// TODO 尽量降低不同模块之间的 耦合
}
