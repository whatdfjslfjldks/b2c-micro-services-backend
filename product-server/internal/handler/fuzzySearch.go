package handler

import (
	"context"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/repository"
)

func (s *Server) FuzzySearch(ctx context.Context, req *pb.FuzzySearchRequest) (
	*pb.FuzzySearchResponse, error) {

	products, totalItems, err := repository.FuzzySearch(req.Keyword, req.CurrentPage, req.PageSize)
	if err != nil {

		return &pb.FuzzySearchResponse{
			Code:       500,
			StatusCode: "GLB-003",
			Msg:        "error",
		}, err
	}
	return &pb.FuzzySearchResponse{
		Code:        200,
		StatusCode:  "GLB-000",
		Msg:         "ok",
		ProductList: products,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
		TotalItems:  totalItems,
	}, err
}
