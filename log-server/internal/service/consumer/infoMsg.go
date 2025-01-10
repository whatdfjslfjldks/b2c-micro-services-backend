package consumer

import (
	"micro-services/log-server/internal/repository"
	"micro-services/log-server/pkg/kafka/model"
)

func InfoMsg(message model.Log) {
	// 将日志存入数据库
	//fmt.Println("在service层info： ", message)
	repository.SaveInfoLogIntoMysql(message)
}
