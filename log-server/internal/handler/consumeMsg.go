package handler

import (
	"context"
	"log"
	c "micro-services/log-server/internal/service/consumer"
	"micro-services/log-server/pkg/kafka"
	"micro-services/log-server/pkg/kafka/model"
)

// ConsumeMsg 用于启动 Kafka 消费者并处理消息
func ConsumeMsg() {
	// 创建 Kafka 消费者实例
	consumer := &kafka.LogKafkaConsumer{}

	// 处理消费的消息，回调函数中的业务逻辑
	err := consumer.ConsumeMessages(context.Background(), func(message model.Log) error {
		// 根据日志级别转发到不同的处理函数
		switch message.Level {
		case "INFO":
			c.InfoMsg(message)
		case "WARN":
			c.WarnMsg(message)
		case "ERROR":
			c.ErrorMsg(message)
		default:
			log.Printf("无法识别的日志级别: %s, 消息内容: %v\n", message.Level, message)
		}
		return nil
	})
	// 如果有错误返回，记录日志
	if err != nil {
		log.Printf("Kafka 消费消息失败: %v", err)
	}
}
