package consumer

import (
	"log"
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/pkg/instance"
	"micro-services/recommend-server/pkg/kafka/model"
)

func PurchaseMsg(message interface{}) {
	msg, ok := message.(model.Recommend)
	if !ok {
		log.Printf("类型断言失败，message 不是 model.Recommend 类型: %v\n", message)
	}

	err := repository.CalAndSaveVectorInToRedis(msg)
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

	err = repository.SavePurchaseMsgIntoMysql(msg)
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
