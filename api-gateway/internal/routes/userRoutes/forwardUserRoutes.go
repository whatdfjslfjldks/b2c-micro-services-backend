package userRoutes

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"micro-services/api-gateway/internal/instance"
	userServerProto "micro-services/pkg/proto/user-server"
)

var model = "user-server"

// SendVerifyCode 发送邮箱验证码
func SendVerifyCode(c *gin.Context) {
	var request userServerProto.EmailSendCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
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
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 使用 proto.Clone 创建响应副本，避免直接复制锁定结构体
	respCopy := proto.Clone(&resp).(*userServerProto.EmailSendCodeResponse)
	c.JSON(int(resp.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// CheckVerifyCode 验证邮箱验证码并登录或注册
func CheckVerifyCode(c *gin.Context) {

	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")

	var request userServerProto.EmailVerifyCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.EmailVerifyCodeRequest{
		Email:      request.Email,
		VerifyCode: request.VerifyCode,
		Ip:         ip,
		UserAgent:  agent,
	}
	var resp userServerProto.EmailVerifyCodeResponse
	err = instance.GrpcClient.CallService(model, "checkVerifyCode", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.EmailVerifyCodeResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"username":     respCopy.Username,
			"userId":       respCopy.UserId,
			"role":         respCopy.Role,
			"accessToken":  respCopy.AccessToken,
			"refreshToken": respCopy.RefreshToken,
			"avatarUrl":    respCopy.Avatar,
		},
	})
}

// LoginByPassword 用户名密码登录
func LoginByPassword(c *gin.Context) {
	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")

	var request userServerProto.UsernameLoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	// TODO 把请求数据再转化一次有必要吗？应该是有的，为了防止接收到多余信息，明确只拿取需要的数据
	req := userServerProto.UsernameLoginRequest{
		Username:  request.Username,
		Password:  request.Password,
		Ip:        ip,
		UserAgent: agent,
	}
	var resp userServerProto.UsernameLoginResponse
	err = instance.GrpcClient.CallService(model, "loginByPassword", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.UsernameLoginResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"username":     respCopy.Username,
			"userId":       respCopy.UserId,
			"role":         respCopy.Role,
			"accessToken":  respCopy.AccessToken,
			"refreshToken": respCopy.RefreshToken,
			"avatarUrl":    respCopy.Avatar,
		},
	})
}

// TestAccessToken 查看访问令牌是否过期
func TestAccessToken(c *gin.Context) {
	var request userServerProto.TestAccessTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.TestAccessTokenRequest{
		AccessToken: request.AccessToken,
	}
	var resp userServerProto.TestAccessTokenResponse
	err = instance.GrpcClient.CallService(model, "testAccessToken", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.TestAccessTokenResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// TestRefreshToken 查看刷新令牌是否过期
func TestRefreshToken(c *gin.Context) {
	var request userServerProto.TestRefreshTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.TestRefreshTokenRequest{
		RefreshToken: request.RefreshToken,
	}
	var resp userServerProto.TestRefreshTokenResponse
	err = instance.GrpcClient.CallService(model, "testRefreshToken", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.TestRefreshTokenResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"accessToken": respCopy.AccessToken,
		},
	})
}

// ChangeUsername 修改用户名
func ChangeUsername(c *gin.Context) {
	var request userServerProto.ChangeUsernameRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.ChangeUsernameRequest{
		UserId:      request.UserId,
		Username:    request.Username,
		AccessToken: request.AccessToken,
	}
	var resp userServerProto.ChangeUsernameResponse
	err = instance.GrpcClient.CallService(model, "changeUsername", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.ChangeUsernameResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"username": respCopy.Username,
		},
	})
}

// ChangeEmail 修改邮箱 TODO 调用之前，先验证旧邮箱，再验证新游戏，此接口只用作存储
func ChangeEmail(c *gin.Context) {
	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")
	var request userServerProto.ChangeEmailRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.ChangeEmailRequest{
		UserId:      request.UserId,
		Email:       request.Email,
		AccessToken: request.AccessToken,
		Ip:          ip,
		UserAgent:   agent,
	}
	var resp userServerProto.ChangeEmailResponse
	err = instance.GrpcClient.CallService(model, "changeEmail", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.ChangeEmailResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// ChangePassword 修改密码，注意判断把密码加密的逻辑加上
func ChangePassword(c *gin.Context) {
	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")
	var request userServerProto.ChangePasswordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.ChangePasswordRequest{
		UserId:      request.UserId,
		NewPassword: request.NewPassword,
		OldPassword: request.OldPassword,
		AccessToken: request.AccessToken,
		Ip:          ip,
		UserAgent:   agent,
	}
	var resp userServerProto.ChangePasswordResponse
	err = instance.GrpcClient.CallService(model, "changePassword", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.ChangePasswordResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// ChangePasswordByEmail 利用邮箱验证码修改密码，密码忘记后重置
func ChangePasswordByEmail(c *gin.Context) {
	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")
	var request userServerProto.ChangePasswordByEmailRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		//fmt.Println("sdfsdf: ", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.ChangePasswordByEmailRequest{
		UserId:      request.UserId,
		Email:       request.Email,
		VerifyCode:  request.VerifyCode,
		NewPassword: request.NewPassword,
		AccessToken: request.AccessToken,
		Ip:          ip,
		UserAgent:   agent,
	}
	var resp userServerProto.ChangePasswordByEmailResponse
	err = instance.GrpcClient.CallService(model, "changePasswordByEmail", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.ChangePasswordByEmailResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// EditUserInfo 修改用户信息
func EditUserInfo(c *gin.Context) {
	var request userServerProto.EditUserInfoRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.EditUserInfoRequest{
		UserId:      request.UserId,
		AvatarUrl:   request.AvatarUrl,
		Bio:         request.Bio,
		Location:    request.Location,
		AccessToken: request.AccessToken,
	}
	var resp userServerProto.EditUserInfoResponse
	err = instance.GrpcClient.CallService(model, "editUserInfo", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*userServerProto.EditUserInfoResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}
