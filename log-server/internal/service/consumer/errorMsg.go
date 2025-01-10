package consumer

import (
	"micro-services/log-server/internal/repository"
	"micro-services/log-server/pkg/kafka/model"
)

func ErrorMsg(message model.Log) {
	repository.SaveErrorLogIntoMysql(message)
}
