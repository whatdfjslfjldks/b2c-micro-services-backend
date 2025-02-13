package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                                                                        // 允许的跨域来源，* 表示允许所有
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                  // 允许的请求方法
		AllowHeaders:     []string{"Access-Token", "Refresh-Token", "X-Real-IP", "X-Forwarded-For", "Origin", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,                                                                                                                 // 是否允许携带凭证（例如 cookies）
		MaxAge:           12 * time.Hour,
	}))

	model := "/api"

	r.POST(model+"/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
			"data": gin.H{
				"user": gin.H{
					"username":    "admin",
					"id":          1,
					"role":        1,
					"status":      1,
					"permissions": 1,
				},
				"accessToken":  "accessToken",
				"refreshToken": "refreshToken",
			},
		})
	})

	log.Println("mgt-server is listening on :8082")
	err := r.Run(":8082")
	if err != nil {
		panic(err)
		return
	}
}
