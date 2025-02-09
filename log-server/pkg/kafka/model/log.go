package model

// Log 定义消息日志格式
type Log struct {
	Source      string `json:"source"`       // 服务名称
	RequestPath string `json:"request_path"` // 请求路径，api接口或是模块通信grpc
	StatusCode  string `json:"status_code"`
	Msg         string `json:"msg"`
	Level       string `json:"level"` // INFO ERROR WARN DEBUG
	Time        string `json:"time"`
}
