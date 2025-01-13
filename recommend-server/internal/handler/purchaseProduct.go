package handler

import (
	"context"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/pkg/utils"
	h "micro-services/recommend-server/pkg/kafka/handler"
	"micro-services/recommend-server/pkg/kafka/model"
)

func (s *Server) PurchaseProduct(ctx context.Context, req *pb.PurchaseProductRequest) (
	*pb.PurchaseProductResponse, error) {
	formattedMsg := model.Recommend{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Status:    "PURCHASE",
		Time:      utils.GetTime(),
	}
	_ = h.KafkaProducer.PublishMessage(formattedMsg, 1)
	return nil, nil
}
