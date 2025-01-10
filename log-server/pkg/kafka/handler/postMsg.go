package handler

import (
	"fmt"
	"log"
	"micro-services/log-server/pkg/kafka/model"
)

func PostMsg(source string,
	requestPath string,
	statusCode string,
	msg string,
	level string,
	time string) {
	//fmt.Println("日志内容： ", source, requestPath, statusCode, msg, level, time)
	// 异步执行，不影响主线程，日志信息允许丢失，容错高
	go func() {
		// 格式化 log
		formattedLog := model.Log{
			Source:      source,
			RequestPath: requestPath,
			StatusCode:  statusCode,
			Msg:         msg,
			Level:       level,
			Time:        time,
		}
		//fmt.Println("日志内容： ", formattedLog)
		// 发布消息 发送到对应分区
		var partition int32
		switch level {
		case "INFO":
			partition = 0
		case "WARN":
			partition = 1
		case "ERROR":
			partition = 2
		default:
			return
		}
		fmt.Println("Sdfasdfsdf:", partition)
		err := KafkaProducer.PublishMessage(formattedLog, partition)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		}
	}()
}
