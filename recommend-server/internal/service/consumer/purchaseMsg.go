package consumer

import (
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/pkg/instance"
	"micro-services/recommend-server/pkg/kafka/model"
)

func PurchaseMsg(message model.Recommend) {
	err := repository.SavePurchaseMsgIntoMysql(message)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/purchaseProduct",
			Source:      "recommend-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}
	return
}
