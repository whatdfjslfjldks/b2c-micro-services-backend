package handler

import (
	"context"
	h "micro-services/log-server/pkg/kafka/handler"
	pb "micro-services/pkg/proto/log-server"
)

// 接受发送日志请求
func (s *Server) PostLog(ctx context.Context, req *pb.PostLogRequest) (
	*pb.PostLogResponse, error) {
	//fmt.Println("receive log", req)
	h.PostMsg(req.Source, req.RequestPath, req.StatusCode, req.Msg, req.Level, req.Time)
	// 不返回值，异步请求，且允许数据丢失
	return nil, nil
}
