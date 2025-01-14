package handler

import (
	"context"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/pkg/utils"
	h "micro-services/recommend-server/pkg/kafka/handler"
	"micro-services/recommend-server/pkg/kafka/model"
)

func (s *Server) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (
	*pb.SearchProductResponse, error) {

	formattedMsg := model.Recommend{
		UserId:  req.UserId,
		Keyword: req.Keyword,
		Status:  "SEARCH",
		Time:    utils.GetTime(),
	}
	_ = h.KafkaProducer.PublishMessage(formattedMsg, 2)
	return nil, nil
}
