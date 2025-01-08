package handler

import (
	"fmt"
	"micro-services/log-server/pkg/kafka"
	"micro-services/log-server/pkg/kafka/model"
)

func ConsumeMsg() {
	consumer := &kafka.LogKafkaConsumer{}
	// 不处理错误
	_ = consumer.ConsumeMessages(func(message model.Log) error {
		// 在这里处理每条消息,转发到service层
		fmt.Println("处理消息: ", message)
		return nil
	})
}
