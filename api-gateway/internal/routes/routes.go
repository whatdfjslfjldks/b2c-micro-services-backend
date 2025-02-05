package routes

import (
	"github.com/gin-gonic/gin"
	"micro-services/api-gateway/internal/routes/aiRoutes"
	"micro-services/api-gateway/internal/routes/payRoutes"
	"micro-services/api-gateway/internal/routes/productRoutes"
	"micro-services/api-gateway/internal/routes/recommendRoutes"
	"micro-services/api-gateway/internal/routes/userRoutes"
)

func SetupRoutes(router *gin.Engine) {
	// -----------------ping-pong-------------------------------
	router.GET("/api/ping", func(c *gin.Context) {
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
		// TODO 有空或者注意到，给多表插入操作加上事件回滚
		productServer.GET("/getProductList", productRoutes.GetProductList)
		productServer.GET("/getSecKillList", productRoutes.GetSecKillList)
		// 获取详情页商品信息 (三种kind一个接口)
		productServer.GET("/getProductById", productRoutes.GetProductById)
		//productServer.GET("/getProductDetailById", productRoutes.GetProductDetailById)
		// TODO： 批量上传接口，后面加一个身份验证，管理员权限才可以，accessToken role=admin
		// 不支持图片的上传，图片可以通过对单一商品的修改上传
		// TODO 返回一个预计上传时间
		// TODO **一定要做身份校验** 因为大文件上传并操控数据库，容易造成数据库雪崩
		productServer.POST("/uploadProductByExcel", productRoutes.UploadProductByExcel)
		// TODO 鉴权
		productServer.POST("/uploadSecKillProduct", productRoutes.UploadSecKillProduct)
	}

	// -----------------处理recommend模块请求--------------------------------
	model3 := "recommend-server"
	recommendServer := router.Group("/api/" + model3)
	{
		// 数据埋点
		// TODO 分开接口，减轻接口在高并发场景下的处理压力
		recommendServer.POST("/clickProduct", recommendRoutes.ClickProduct)
		recommendServer.POST("/browseProduct", recommendRoutes.BrowseProduct)
		recommendServer.POST("/purchaseProduct", recommendRoutes.PurchaseProduct)
		recommendServer.POST("/searchProduct", recommendRoutes.SearchProduct)

		// 获取推荐商品
		recommendServer.GET("/GetRecommendProductList", recommendRoutes.GetRecommendProductList)

	}

	// -----------------处理ai模块请求--------------------------------
	model4 := "ai-server"
	aiServer := router.Group("/api/" + model4)
	{
		aiServer.POST("/talk", aiRoutes.Talk)
	}

	// -----------------处理pay模块请求--------------------------------
	model5 := "pay-server"
	payServer := router.Group("/api/" + model5)
	{
		payServer.POST("/tradePreCreate", payRoutes.TradePreCreate)
	}

}
