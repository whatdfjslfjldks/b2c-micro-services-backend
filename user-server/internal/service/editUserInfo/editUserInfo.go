package editUserInfo

import (
	pb "micro-services/pkg/proto/user-server"
	"micro-services/user-server/pkg/config"
	"time"
)

func EditUserInfo(userId int64, avatarUrl string, bio string, location string) *pb.EditUserInfoResponse {
	currentTime := time.Now()
	query := "UPDATE b2c_user.user_profiles SET avatar_url=?,bio=?,location=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, avatarUrl, bio, location, currentTime, userId)
	if err != nil {
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
