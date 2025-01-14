package consumer

import (
	"log"
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/pkg/instance"
	"micro-services/recommend-server/pkg/kafka/model"
)

func ClickMsg(message interface{}) {
	msg, ok := message.(model.Recommend)
	if !ok {
		log.Printf("类型断言失败，message 不是 model.Recommend 类型: %v\n", message)
	}
	// TODO 这里进行增量计算
	//fmt.Println("点击消息：", msg)
	err := repository.CalAndSaveVectorInToRedis(msg)
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

	err = repository.SaveClickMsgIntoMysql(msg)
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
