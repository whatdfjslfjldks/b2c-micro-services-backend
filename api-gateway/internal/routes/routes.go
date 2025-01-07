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
		userServer.POST("/testAccessToken", userRoutes.TestAccessToken)
		userServer.POST("/testRefreshToken", userRoutes.TestRefreshToken)
		userServer.POST("/changeUsername", userRoutes.ChangeUsername)
		userServer.POST("/changeEmail", userRoutes.ChangeEmail)
		userServer.POST("/changePassword", userRoutes.ChangePassword)
		// TODO 风控模块---增加一个判断异常，短时间多次修改密码或其他，提醒用户异常
		userServer.POST("/changePasswordByEmail", userRoutes.ChangePasswordByEmail)
		userServer.POST("/editUserInfo", userRoutes.EditUserInfo)

	}
}
