package repository

import (
	"log"
	"micro-services/log-server/pkg/config"
	"micro-services/log-server/pkg/kafka/model"
)

func SaveInfoLogIntoMysql(message model.Log) {
	query := "INSERT INTO b2c_log_info.logs (source, request_path, status_code, msg, level, time) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.MySqlInfoClient.Exec(query, message.Source, message.RequestPath, message.StatusCode, message.Msg, message.Level, message.Time)
	if err != nil {
		log.Printf("Failed to insert into MySQL: %v", err)
	}
}
func SaveWarnLogIntoMysql(message model.Log) {
	query := "INSERT INTO b2c_log_warn.logs (source, request_path, status_code, msg, level, time) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.MySqlWarnClient.Exec(query, message.Source, message.RequestPath, message.StatusCode, message.Msg, message.Level, message.Time)
	if err != nil {
		log.Printf("Failed to insertinto MySQL: %v", err)
	}
}
func SaveErrorLogIntoMysql(message model.Log) {
	query := "INSERT INTO b2c_log_error.logs (source, request_path, status_code, msg, level, time) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.MySqlErrorClient.Exec(query, message.Source, message.RequestPath, message.StatusCode, message.Msg, message.Level, message.Time)
	if err != nil {
		log.Printf("Failed to insert into MySQL: %v", err)
	}
}
