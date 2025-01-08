package handler

import (
	pb "micro-services/pkg/proto/log-server"
)

type Server struct {
	pb.UnimplementedLogServiceServer
}
