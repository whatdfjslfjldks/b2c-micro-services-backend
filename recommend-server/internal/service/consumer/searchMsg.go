package consumer

import (
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/pkg/instance"
	"micro-services/recommend-server/pkg/kafka/model"
)

func SearchMsg(message model.Recommend) {
	err := repository.SaveSearchMsgIntoMysql(message)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/searchProduct",
			Source:      "recommend-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}
}
