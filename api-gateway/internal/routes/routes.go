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
		// TODO 检测非常用ip和agent
		userServer.POST("/checkVerifyCode", userRoutes.CheckVerifyCode)
		// TODO 检测非常用ip和agent
		userServer.POST("/loginByPassword", userRoutes.LoginByPassword)
		userServer.POST("/testAccessToken", userRoutes.TestAccessToken)
		userServer.POST("/testRefreshToken", userRoutes.TestRefreshToken)
		userServer.POST("/changeUsername", userRoutes.ChangeUsername)
		// TODO 检测非常用ip和agent
		userServer.POST("/changeEmail", userRoutes.ChangeEmail)
		// TODO 检测非常用ip和agent
		userServer.POST("/changePassword", userRoutes.ChangePassword)
		// TODO 检测非常用ip和agent
		userServer.POST("/changePasswordByEmail", userRoutes.ChangePasswordByEmail)
		userServer.POST("/editUserInfo", userRoutes.EditUserInfo)

	}
}
