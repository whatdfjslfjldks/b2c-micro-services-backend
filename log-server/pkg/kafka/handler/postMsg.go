package handler

import (
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
		// 发布消息
		err := KafkaProducer.PublishMessage(formattedLog)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		}
	}()
}
