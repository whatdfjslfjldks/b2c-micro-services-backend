package handler

import (
	"context"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/pkg/utils"
	h "micro-services/recommend-server/pkg/kafka/handler"
	"micro-services/recommend-server/pkg/kafka/model"
)

func (s *Server) BrowseProduct(ctx context.Context, req *pb.BrowseProductRequest) (
	*pb.BrowseProductResponse, error) {
	formattedMsg := model.Recommend{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Status:    "BROWSE",
		Time:      utils.GetTime(),
	}
	_ = h.KafkaProducer.PublishMessage(formattedMsg, 2)
	return nil, nil
}
