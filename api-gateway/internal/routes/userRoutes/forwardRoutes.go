package userRoutes

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"micro-services/api-gateway/internal/instance"
	userServerProto "micro-services/pkg/proto/user-server"
	"net/http"
)

var model = "user-server"

// 发送邮箱验证码
func SendVerifyCode(c *gin.Context) {
	var request userServerProto.EmailSendCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		return
	}
	//构建请求对象
	req := userServerProto.EmailSendCodeRequest{
		Email: request.Email,
	}
	//创建响应对象
	var resp userServerProto.EmailSendCodeResponse
	err = instance.GrpcClient.CallService(model, "sendVerifyCode", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "服务器错误: " + err.Error(),
		})
		return
	}
	// 使用 proto.Clone 创建响应副本，避免直接复制锁定结构体
	respCopy := proto.Clone(&resp).(*userServerProto.EmailSendCodeResponse)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  respCopy.Msg,
	})
}

// 验证邮箱验证码并登录或注册
func CheckVerifyCode(c *gin.Context) {
	var request userServerProto.EmailVerifyCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "无法解析请求: " + err.Error(),
		})
		return
	}
	req := userServerProto.EmailVerifyCodeRequest{
		Email:      request.Email,
		VerifyCode: request.VerifyCode,
	}
	var resp userServerProto.EmailVerifyCodeResponse
	err = instance.GrpcClient.CallService(model, "checkVerifyCode", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "服务器错误：" + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.EmailVerifyCodeResponse)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "响应成功",
		"data": respCopy,
	})
}

// 用户名密码登录
func LoginByPassword(c *gin.Context) {
	var request userServerProto.UsernameLoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "无法解析请求: " + err.Error(),
		})
		return
	}
	// TODO 把请求数据再转化一次有必要吗？应该是有的，为了防止接收到多余信息，明确只拿取需要的数据
	req := userServerProto.UsernameLoginRequest{
		Username: request.Username,
		Password: request.Password,
	}
	var resp userServerProto.UsernameLoginResponse
	err = instance.GrpcClient.CallService(model, "loginByPassword", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "服务器错误：" + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.UsernameLoginResponse)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "响应成功",
		"data": respCopy,
	})
}
