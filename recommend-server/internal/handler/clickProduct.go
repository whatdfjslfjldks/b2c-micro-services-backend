package handler

import (
	"context"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/pkg/utils"
	h "micro-services/recommend-server/pkg/kafka/handler"
	"micro-services/recommend-server/pkg/kafka/model"
)

func (s *Server) ClickProduct(ctx context.Context, req *pb.ClickProductRequest) (
	*pb.ClickProductResponse, error) {
	//存入 kafka
	formattedMsg := model.Recommend{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Status:    "CLICK",
		Time:      utils.GetTime(),
	}
	_ = h.KafkaProducer.PublishMessage(formattedMsg, 0)
	return nil, nil
}
