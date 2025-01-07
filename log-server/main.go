package main

import (
	"fmt"
	"micro-services/pkg/kafka/logServerKafka"
)

func main() {
	// 测试
	fmt.Println("begin----------------------")
	// 初始化kafka客户端
	err := logServerKafka.InitKafkaConfig()
	if err != nil {
		fmt.Println("init kafka config failed")
		return
	}

}
