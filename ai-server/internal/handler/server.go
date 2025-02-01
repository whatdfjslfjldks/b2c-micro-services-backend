package handler

import (
	pb "micro-services/pkg/proto/ai-server"
)

type Server struct {
	pb.UnimplementedAIServiceServer
}
