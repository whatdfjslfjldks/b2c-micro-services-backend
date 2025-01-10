package riskFilter

//
//import (
//	"micro-services/api-gateway/internal/instance"
//	pb "micro-services/pkg/proto/risk-server"
//)
//
//func RiskIpAndAgent(ip string, agent string) (
//	*pb.RiskLoginIpAndAgentResponse, error) {
//
//	//fmt.Println("ip:", ip, "agent:", agent)
//	// 转发到风控模块
//	err := instance.GrpcClient.CallService("risk-server", "RiskLoginIpAndAgent")
//	return nil, nil
//}

// TODO 风控拦截好像不能再这里，应该在服务内部使用，这里可以做ip白名单，黑名单处理
