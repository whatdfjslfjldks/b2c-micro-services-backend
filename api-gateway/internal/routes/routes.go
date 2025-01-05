package routes

import (
	"github.com/gin-gonic/gin"
	"micro-services/api-gateway/internal/routes/userRoutes"
)

func SetupRoutes(router *gin.Engine) {
	// -----------------ping-pong-------------------------------
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// -----------------处理user模块请求--------------------------------
	model := "user-server"
	userServer := router.Group("/api/" + model)
	{
		userServer.POST("/sendVerifyCode", userRoutes.SendVerifyCode)
		userServer.POST("/checkVerifyCode", userRoutes.CheckVerifyCode)
		userServer.POST("/loginByPassword", userRoutes.LoginByPassword)

	}
}
