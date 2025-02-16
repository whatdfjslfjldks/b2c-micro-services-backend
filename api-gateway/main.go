package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api-gateway/internal/instance"
	"micro-services/api-gateway/internal/routes"
	"micro-services/api-gateway/pkg/localLog"
	"time"
)

// TODO: etcd注册中心  50060-50063：user-server
// TODO:             50052：log-server
// TODO:             50053：risk-server
// TODO：            50054：product-server
// TODO:             50055：recommend-server
// TODO:             50056：ai-server
// TODO:             50057：pay-server
// TODO:             50058: order-server
func main() {
	err := localLog.InitLog()
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}
	// 启动 HTTP 服务
	r := gin.Default()
	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                                                                        // 允许的跨域来源，* 表示允许所有
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                  // 允许的请求方法
		AllowHeaders:     []string{"Access-Token", "Refresh-Token", "X-Real-IP", "X-Forwarded-For", "Origin", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,                                                                                                                 // 是否允许携带凭证（例如 cookies）
		MaxAge:           12 * time.Hour,                                                                                                       // 预检请求的有效期，单位是时间
	}))
	r.Use(TimeMW())

	// 用gin做网关，进行路由的接收和转发
	routes.SetupRoutes(r)

	//创建实例
	instance.NewInstance()

	localLog.GateWayLog.Info("api-gateway start: HTTP server is listening on port 8080")
	// 启动 HTTP 服务，监听 8080 端口
	log.Println("HTTP server is listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	// TODO: 风控模块，日志模块，商品模块，订单模块，支付模块
	// TODO: 尽量降低不同模块之间的 耦合
	// TODO: 风控和日志模块待更新，现在只做了一个示范
	// TODO: 数据库初始化封装
}

// TimeMW 响应耗时中间件
func TimeMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		ctx.Next()
		timeCost := time.Since(now)
		// TODO: 信息收集
		_ = fmt.Sprintf("request %s cost time %d ms\n", ctx.Request.URL.Path, timeCost.Milliseconds())
		//fmt.Println("消息： ", msg)
		//localLog.GateWayLog.Info(msg)
	}
}
