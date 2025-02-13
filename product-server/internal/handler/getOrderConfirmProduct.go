package handler

import (
	"context"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/repository"
)

func (s *Server) GetOrderConfirmProduct(ctx context.Context, req *pb.GetOrderConfirmProductRequest) (
	*pb.GetOrderConfirmProductResponse, error) {
	//fmt.Println("111", req.ProductId)
	r, err := repository.GetOrderConfirmProduct(req.ProductId)
	if err != nil {
		log.Printf("GetOrderConfirmProduct error: %v", err)
		a := &pb.GetOrderConfirmProductResponse{
			Code:       500,
			Msg:        "GetOrderConfirmProduct error",
			StatusCode: "GLB-003",
		}
		return a, nil
	}

	return &pb.GetOrderConfirmProductResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "获取数据成功！",
		Products:   r,
	}, nil
}
