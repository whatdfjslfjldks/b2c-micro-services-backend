package emailService

import (
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/repository"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/instance"
	"micro-services/user-server/pkg/token"
)

// TODO 尽量不要再service层写数据库操作，保证每层分工明确
func VerifyCode(email string, code string) bool {
	err := repository.GetCode(email, code)
	if err != nil {
		return false
	}
	//fmt.Println(result)
	return true
}

// LoginByEmail TODO 多次操作数据库，非常耗时，待优化！
func LoginByEmail(email string) (resp *pb.EmailVerifyCodeResponse, err error) {
	resp = &pb.EmailVerifyCodeResponse{}
	result := repository.IsEmailExist(email)
	if result {
		//邮箱存在
		//查找基本信息
		id, name, role, err := repository.GetUserInfoByEmail(email)
		if err != nil {
			resp.Code = 400
			resp.StatusCode = "USR-002"
			resp.Msg = "用户不存在！"
			return resp, nil
		}
		if id != nil {
			resp.UserId = *id
		}
		resp.Username = name
		resp.Role = role
		//查找头像
		avatarUrl, err := repository.GetAvatarUrlById(resp.UserId)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "操作数据库出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		resp.Avatar = avatarUrl
		//生成双token
		refreshToken, err := token.GenerateRefreshToken(resp.UserId, resp.Role)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "USR-003"
			resp.Msg = "refreshToken 生成出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "USR-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		accessToken, err := token.GenerateAccessToken(resp.UserId, resp.Role)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "USR-003"
			resp.Msg = "accessToken 生成出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "USR-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		resp.RefreshToken = refreshToken
		resp.AccessToken = accessToken
		//把双token存入redis数据库
		err = repository.SaveToken(resp.UserId, resp.RefreshToken, resp.AccessToken)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "操作数据库出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		resp.Code = 200
		resp.StatusCode = "GLB-000"
		resp.Msg = "登录成功！"
		return resp, nil
	} else {
		//fmt.Println("邮箱不不不不不存在的情况----------------")
		//邮箱不存在
		resp.Username = userPkg.GenerateUsername()
		resp.Role = "user"
		userId, err := repository.SaveUserInfo(resp.Username, email, resp.Role)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "操作数据库出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		resp.UserId = userId
		//生成双token
		refreshToken, err := token.GenerateRefreshToken(resp.UserId, resp.Role)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "USR-003"
			resp.Msg = "refreshToken 生成出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "USR-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		accessToken, err := token.GenerateAccessToken(resp.UserId, resp.Role)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "USR-003"
			resp.Msg = "accessToken 生成出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "USR-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		resp.RefreshToken = refreshToken
		resp.AccessToken = accessToken
		//把双token存入redis数据库
		err = repository.SaveToken(resp.UserId, resp.RefreshToken, resp.AccessToken)
		if err != nil {
			resp.Code = 500
			resp.StatusCode = "GLB-003"
			resp.Msg = "操作数据库出错！"
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/loginByEmail",
				Source:      "user-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return resp, nil
		}
		// 头像置为默认值，后面可以使用服务器上的默认图片
		resp.Avatar = "default"

		return resp, nil
	}

}
