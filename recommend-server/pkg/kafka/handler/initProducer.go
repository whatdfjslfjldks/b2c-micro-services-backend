package handler

import (
	"fmt"
	"micro-services/recommend-server/pkg/kafka"
)

var KafkaProducer *kafka.RecommendKafkaProducer

func InitProducer() {
	// 初始化 Kafka 配置
	err := kafka.InitKafkaConfig()
	if err != nil {
		fmt.Println("init kafka config failed:", err)
		return
	}

	// 创建生产者实例
	KafkaProducer = &kafka.RecommendKafkaProducer{}

	// 初始化 Kafka 生产者
	err = KafkaProducer.InitProducer()
	if err != nil {
		fmt.Println("init kafka producer failed:", err)
		return
	}

}
