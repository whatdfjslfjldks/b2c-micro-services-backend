package handler

import (
	"context"
	"fmt"
	pb "micro-services/pkg/proto/product-server"
)

func (s *Server) GetProductList(ctx context.Context, req *pb.GetProductListRequest) (
	*pb.GetProductListResponse, error) {
	//
	fmt.Println("sdfds", req.CurrentPage, req.PageSize, req.CategoryId)

	return nil, nil
}
