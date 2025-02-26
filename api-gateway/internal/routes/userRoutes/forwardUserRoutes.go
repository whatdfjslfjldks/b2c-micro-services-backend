package userRoutes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api-gateway/internal/instance"
	userServerProto "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "sendVerifyCode", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.EmailSendCodeResponse)
	c.JSON(int(respCopy.Code), gin.H{
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "checkVerifyCode", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.EmailVerifyCodeResponse)
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
			"email":        respCopy.Email,
			"bio":          respCopy.Bio,
			"createAt":     respCopy.CreateAt,
		},
	})
}

// LoginByPassword 用户名密码登录
func LoginByPassword(c *gin.Context) {
	fmt.Println("12312")
	// 获取请求头里的ip和agent
	ip := c.ClientIP()
	agent := c.Request.Header.Get("User-Agent")
	fmt.Println("1111: ", ip, agent)

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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "loginByPassword", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.UsernameLoginResponse)
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
			"email":        respCopy.Email,
			"bio":          respCopy.Bio,
			"createAt":     respCopy.CreateAt,
		},
	})
}

// TestAccessToken 查看访问令牌是否过期
func TestAccessToken(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	req := userServerProto.TestAccessTokenRequest{
		AccessToken: accessToken,
	}
	var resp interface{}
	err := instance.GrpcClient.CallService(model, "testAccessToken", &req, &resp)
	if err != nil {
		c.String(400, "nook")
		return
		//c.JSON(500, gin.H{
		//	"code":        500,
		//	"status_code": "GLB-002",
		//	"msg":         "grpc调用错误: " + err.Error(),
		//})
		//return
	}
	respCopy := (resp).(*userServerProto.TestAccessTokenResponse)
	if respCopy.Code == 200 {
		c.String(200, "ok")
		return
	} else {
		c.String(400, "nook")
		return
	}
	//c.JSON(int(respCopy.Code), gin.H{
	//	"code":        respCopy.Code,
	//	"status_code": respCopy.StatusCode,
	//	"msg":         respCopy.Msg,
	//})
}

// TestRefreshToken 查看刷新令牌是否过期
func TestRefreshToken(c *gin.Context) {
	refreshToken := c.Request.Header.Get("Refresh-Token")
	req := userServerProto.TestRefreshTokenRequest{
		RefreshToken: refreshToken,
	}
	var resp interface{}
	err := instance.GrpcClient.CallService(model, "testRefreshToken", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.TestRefreshTokenResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"accessToken":  respCopy.AccessToken,
			"refreshToken": respCopy.RefreshToken,
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "changeUsername", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.ChangeUsernameResponse)
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "changeEmail", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.ChangeEmailResponse)
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "changePassword", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.ChangePasswordResponse)
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "changePasswordByEmail", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.ChangePasswordByEmailResponse)
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
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "editUserInfo", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.EditUserInfoResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

func GetUserInfoByUserId(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	req := userServerProto.GetUserInfoByUserIdRequest{
		AccessToken: accessToken,
	}
	var resp interface{}
	err := instance.GrpcClient.CallService(model, "getUserInfoByUserId", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.GetUserInfoByUserIdResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"avatar_url": respCopy.AvatarUrl,
			"name":       respCopy.Name,
			"email":      respCopy.Email,
			"user_id":    respCopy.UserId,
			"bio":        respCopy.Bio,
			"create_at":  respCopy.CreateAt,
		},
	})
}

func UploadAvatar(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	fileContent, e := utils.ReadFileContent(file)
	if e != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "PRT-001",
			"msg":         "读取 图片 文件内容错误！",
		})
		return
	}
	req := userServerProto.UploadAvatarRequest{
		AccessToken: accessToken,
		File:        fileContent,
	}
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "uploadAvatar", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.UploadAvatarResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"avatar_url": respCopy.AvatarUrl,
		},
	})
}

func UpdateName(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	type a struct {
		Name string `json:"name"`
	}
	var b a
	err := c.ShouldBindJSON(&b)
	if err != nil {
		log.Printf("err: %v", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.UpdateNameRequest{
		AccessToken: accessToken,
		Name:        b.Name,
	}
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "updateName", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.UpdateNameResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"name": respCopy.Name,
		},
	})
}
func UpdateBio(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	type a struct {
		Bio string `json:"bio"`
	}
	var b a
	err := c.ShouldBindJSON(&b)
	if err != nil {
		log.Printf("err: %v", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := userServerProto.UpdateBioRequest{
		AccessToken: accessToken,
		Bio:         b.Bio,
	}
	var resp interface{}
	err = instance.GrpcClient.CallService(model, "updateBio", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*userServerProto.UpdateBioResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"bio": respCopy.Bio,
		},
	})
}
