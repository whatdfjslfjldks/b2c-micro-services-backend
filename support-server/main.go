package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/support-server/internal"
	"micro-services/support-server/pkg/config"
	"time"
)

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

// 客服模块，不用grpc，单独拆分出来
func main() {

	initConfig()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                                                                        // 允许的跨域来源，* 表示允许所有
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                  // 允许的请求方法
		AllowHeaders:     []string{"Access-Token", "Refresh-Token", "X-Real-IP", "X-Forwarded-For", "Origin", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,                                                                                                                 // 是否允许携带凭证（例如 cookies）
		MaxAge:           12 * time.Hour,
	}))

	//http.HandleFunc("/chat", handleChat)       // 用户连接
	//http.HandleFunc("/support", handleSupport) // 客服连接
	model := "/api/support-server"

	// 客服注册
	r.POST(model+"/register", internal.SupportRegister)
	// 客服登录
	r.POST(model+"/login", internal.SupportLogin)
	// 用户查找是否有客服在线,如果有就随机选择一个客服返回，后续可优化，“负载均衡”
	r.POST(model+"/findSupport", internal.FindSupport)

	// 建立客服与服务端ws连接，用于推送用户连接消息
	r.GET(model+"/supportConnect", internal.SupportConnect)

	// 连接用户到客服，并生成房间号
	r.GET(model+"/connect", internal.Connect)

	// 连接客服到用户生成的房间号
	r.GET(model+"/connectRoom", internal.ConnectRoom)

	log.Println("Server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
