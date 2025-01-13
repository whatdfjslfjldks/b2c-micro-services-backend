package handler

import (
	"context"
	pb "micro-services/pkg/proto/recommend-server"
)

// TODO ***注意*** 这个需要什么根据keyword判断是哪些id

func (s *Server) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (
	*pb.SearchProductResponse, error) {

	return nil, nil
}
