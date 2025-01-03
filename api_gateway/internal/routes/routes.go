package routes

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"micro-services/api_gateway/internal/grpcClient"
	"micro-services/pkg/etcd"
	user_server_proto "micro-services/pkg/proto/user_server"
	"net/http"
	"time"
)

func SetupRoutes(router *gin.Engine) {
	// -----------------ping-pong-------------------------------
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// -----------------创建etcd客户端--------------------------------
	etcdClient, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		// 丸辣！🌶
		panic(err)
	}
	//defer etcdClient.Close()
	// -----------------获取etcd所有服务的地址,适合微服务变动不频繁，小型项目--------------------------------
	//services, err := etcdClient.GetAllServices()
	// -----------------创建GRPCClient实例--------------------------------
	grpcClient := user_server_grpcClient.NewGRPCClient(etcdClient)
	// -----------------处理user模块请求--------------------------------
	model := "user_server"
	userServer := router.Group("/api/" + model)
	{
		userServer.POST("/sendVerifyCode", func(c *gin.Context) {

			var request user_server_proto.EmailSendCodeRequest
			err := c.ShouldBindJSON(&request)
			if err != nil {
				return
			}
			//构建请求对象
			req := user_server_proto.EmailSendCodeRequest{
				Email: request.Email,
			}
			//创建响应对象
			var resp user_server_proto.EmailSendCodeResponse
			err = grpcClient.CallService(model, "sendVerifyCode", &req, &resp)
			if err != nil {
				c.JSON(500, gin.H{
					"message": "failed to call UserRegister",
				})
				return
			}
			// 使用 proto.Clone 创建响应副本，避免直接复制锁定结构体
			respCopy := proto.Clone(&resp).(*user_server_proto.EmailSendCodeResponse)
			c.JSON(http.StatusOK, respCopy)
		})
	}
}
