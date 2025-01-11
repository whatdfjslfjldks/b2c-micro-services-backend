package handler

import pb "micro-services/pkg/proto/recommend-server"

type Server struct {
	pb.UnimplementedRecommendServiceServer
}
