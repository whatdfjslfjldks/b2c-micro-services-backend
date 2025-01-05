package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api-gateway/internal/instance"
	"micro-services/api-gateway/internal/routes"
)

func main() {
	// 启动 HTTP 服务s
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
}
