package routes

import (
	"github.com/gin-gonic/gin"
	"micro-services/api-gateway/internal/routes/productRoutes"
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

	// -----------------处理product模块请求--------------------------------
	model2 := "product-server"
	productServer := router.Group("/api/" + model2)
	{
		productServer.GET("/getProductList", productRoutes.GetProductList)

		// TODO： 批量上传接口，后面加一个身份验证，管理员权限才可以，accessToken role=admin
		// 不支持图片的上传，图片可以通过对单一商品的修改上传
		// TODO 返回一个预计上传时间
		// TODO **一定要做身份校验** 因为大文件上传并操控数据库，容易造成数据库雪崩
		productServer.POST("/uploadProductByExcel", productRoutes.UploadProductByExcel)
	}

}
