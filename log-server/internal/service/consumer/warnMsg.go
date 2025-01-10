package consumer

import (
	"micro-services/log-server/internal/repository"
	"micro-services/log-server/pkg/kafka/model"
)

func WarnMsg(message model.Log) {
	repository.SaveWarnLogIntoMysql(message)
}
