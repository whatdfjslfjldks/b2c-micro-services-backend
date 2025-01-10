package editUserInfo

import (
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/pkg/config"
	"micro-services/user-server/pkg/instance"
)

func EditUserInfo(userId int64, avatarUrl string, bio string, location string) *pb.EditUserInfoResponse {
	currentTime := utils.GetTime()
	query := "UPDATE b2c_user.user_profiles SET avatar_url=?,bio=?,location=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, avatarUrl, bio, location, currentTime, userId)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/editUserInfo",
			Source:      "user-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		return &pb.EditUserInfoResponse{
			Code:       500,
			StatusCode: "GLB-003",
			Msg:        "数据库错误!",
		}
	} else {
		return &pb.EditUserInfoResponse{
			Code:       200,
			StatusCode: "GLB-000",
			Msg:        "信息修改成功!",
		}
	}

}
