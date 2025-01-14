package handler

import (
	"context"
	"log"
	c "micro-services/recommend-server/internal/service/consumer"
	"micro-services/recommend-server/pkg/kafka"
	"micro-services/recommend-server/pkg/kafka/model"
)

// ConsumeMsg 用于启动 Kafka 消费者并处理消息
func ConsumeMsg() {
	// 创建 Kafka 消费者实例
	consumer := &kafka.RecommendKafkaConsumer{}

	err := consumer.ConsumeMessages(context.Background(), func(message model.Recommend) error {

		//fmt.Println("hello : ", message)
		//return nil

		// 根据日志级别转发到不同的处理函数
		switch message.Status {
		case "CLICK":
			c.ClickMsg(message)
		case "PURCHASE":
			c.PurchaseMsg(message)
		case "BROWSE":
			c.BrowseMsg(message)
			// TODO  search
		case "SEARCH":
			c.SearchMsg(message)
		default:
			log.Printf("无法识别的日志级别: %s, 消息内容: %v\n", message.Status, message)
		}

		return nil
	})

	// 如果有错误返回，记录日志
	if err != nil {
		log.Printf("Kafka 消费消息失败: %v", err)
	}
}
