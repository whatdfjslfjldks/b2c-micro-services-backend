package consumer

import (
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/pkg/instance"
	"micro-services/recommend-server/pkg/kafka/model"
)

func ClickMsg(message model.Recommend) {
	err := repository.SaveClickMsgIntoMysql(message)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/clickProduct",
			Source:      "recommend-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}
	return
}
